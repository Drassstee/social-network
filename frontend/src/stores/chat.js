import { defineStore } from 'pinia'
import { api } from '../api/api'
import { useNotificationStore } from './notifications'

export const useChatStore = defineStore('chat', {
  state: () => ({
    socket: null,
    onlineUsers: [],
    messages: {}, // userID -> messages[]
    groupMessages: {}, // groupID -> messages[]
    unreadCounts: {}, // identifier (u_ID or g_ID) -> count
    activeChatUser: null,
    activeGroupID: null,
    connected: false,
  }),
  actions: {
    connect() {
      if (this.socket) return

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = window.location.host
      this.socket = new WebSocket(`${protocol}//${host}/api/v1/ws`)

      this.socket.onopen = () => {
        this.connected = true
        this.fetchChatableUsers()
        this.fetchUnreadCounts()
      }

      this.socket.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data)
          this.handleSocketMessage(msg)
        } catch (e) {
          console.error('Failed to parse socket message:', e)
        }
      }

      this.socket.onclose = () => {
        this.connected = false
        this.socket = null
        if (this.isUserLoggedIn()) {
          setTimeout(() => this.connect(), 3000)
        }
      }
    },
    disconnect() {
      if (this.socket) {
        this.socket.onclose = null
        this.socket.close()
        this.socket = null
        this.connected = false
      }
    },
    isUserLoggedIn() {
      return !!localStorage.getItem('user')
    },
    async fetchChatableUsers() {
      try {
        const data = await api.get('/chat/users')
        this.onlineUsers = data // This actually contains all chatable users enriched with online status
      } catch (err) {
        console.error('Failed to fetch chatable users:', err)
      }
    },
    async fetchUnreadCounts() {
        try {
            // Private counts
            const pCounts = await api.get('/chat/unread')
            for (const [id, count] of Object.entries(pCounts)) {
                this.unreadCounts[`u_${id}`] = count
            }

            // Group counts
            const gCounts = await api.get('/groups/unread')
            if (Array.isArray(gCounts)) {
                gCounts.forEach(c => {
                    if (c.unread_count > 0) {
                        this.unreadCounts[`g_${c.group_id}`] = c.unread_count
                    }
                })
            }
        } catch (err) {
            console.error('Failed to fetch unread counts:', err)
        }
    },
    handleSocketMessage(msg) {
      switch (msg.type) {
        case 'status_update':
          this.updateUserStatus(msg.data)
          break
        case 'private_message':
          this.addMessage(msg.data)
          break
        case 'group_message':
          this.addGroupMessage(msg.data)
          break
        case 'notification': {
          const notifStore = useNotificationStore()
          notifStore.addNotification(msg.data)
          break
        }
      }

    },
    updateUserStatus(data) {
      const user = this.onlineUsers.find(u => u.id === data.user_id)
      if (user) {
        user.is_online = data.online
      } else {
        // If user not in list, we might need to refresh or just add if we want to show non-followers who just came online (optional)
        // For now, let's just refresh if it's someone we don't know but we might care about
        // Or better, just add them if online
        if (data.online) {
            this.onlineUsers.push({ id: data.user_id, username: data.username, is_online: true })
        }
      }
    },
    // Internal helper for deduplicated message insertion
    _addMessageToStore(targetMap, id, msg) {
      if (!targetMap[id]) {
        targetMap[id] = []
      }
      if (!targetMap[id].some(m => m.id === msg.id)) {
        targetMap[id].push(msg)
      }
    },
    async addMessage(msg) {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      const myID = authStore.user?.id
      
      const otherID = msg.sender_id === myID ? msg.receiver_id : msg.sender_id
      this._addMessageToStore(this.messages, otherID, msg)

      // Increment unread if not active chat
      if (msg.sender_id !== myID && (!this.activeChatUser || this.activeChatUser.id !== otherID)) {
        const key = `u_${otherID}`
        this.unreadCounts[key] = (this.unreadCounts[key] || 0) + 1
        
        // Global toast notification
        const { useUIStore } = await import('./ui')
        useUIStore().showToast(`NEW_MESSAGE: ${msg.body.substring(0, 20)}${msg.body.length > 20 ? '...' : ''}`, 'info')
      } else if (msg.sender_id !== myID) {
          this.markPrivateAsRead(otherID)
      }
    },
    async markPrivateAsRead(senderID) {
      try {
        await api.post(`/chat/read?sender_id=${senderID}`)
        delete this.unreadCounts[`u_${senderID}`]
      } catch (err) {
        console.error('Failed to mark private as read:', err)
      }
    },
    async uploadImage(file) {
      const formData = new FormData()
      formData.append('image', file)
      try {
        const data = await api.post('/chat/upload', formData)
        return data.url
      } catch (err) {
        console.error('Chat image upload failed:', err)
        throw err
      }
    },
    async sendPrivateMessage(receiverID, body, file = null) {
      let imageURL = null
      if (file) {
        try {
          imageURL = await this.uploadImage(file)
        } catch (err) {
          console.error('Image upload failed, message not sent', err)
          throw err
        }
      }

      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
        throw new Error('COMM_LINK_OFFLINE')
      }

      const payload = {
        type: 'private_message',
        data: {
          receiver_id: receiverID,
          body: body,
          image_url: imageURL
        }
      }
      this.socket.send(JSON.stringify(payload))
    },
    async fetchMessages(userID) {
      try {
        const data = await api.get(`/chat/messages?user_id=${userID}&limit=50&offset=0`)
        // Reverse because backend returns newest first (DESC), but we display oldest to newest
        this.messages[userID] = (data || []).reverse()
        
        // Clear unread
        this.markPrivateAsRead(userID)
      } catch (err) {
        console.error('Failed to fetch messages:', err)
      }
    },
    async addGroupMessage(msg) {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      const myID = authStore.user?.id
      
      const groupID = msg.group_id
      this._addMessageToStore(this.groupMessages, groupID, msg)

      // Increment unread if not active group
      if (Number(msg.sender_id) !== Number(myID) && Number(this.activeGroupID) !== Number(groupID)) {
        const key = `g_${groupID}`
        this.unreadCounts[key] = (this.unreadCounts[key] || 0) + 1

        // Global toast notification
        const { useUIStore } = await import('./ui')
        useUIStore().showToast(`GROUP_MSG: ${msg.body.substring(0, 20)}${msg.body.length > 20 ? '...' : ''}`, 'info')
      } else if (Number(msg.sender_id) !== Number(myID)) {
          this.markGroupAsRead(groupID)
      }
      return msg
    },
    async markGroupAsRead(groupID) {
      try {
        await api.post(`/groups/${groupID}/read`)
        delete this.unreadCounts[`g_${groupID}`]
      } catch (err) {
        console.error('Failed to mark group as read:', err)
      }
    },
    async fetchGroupMessages(groupID) {
      try {
        const data = await api.get(`/groups/${groupID}/messages`)
        // Reverse because backend returns newest first (DESC), but we display oldest to newest
        this.groupMessages[groupID] = (data || []).reverse()
        
        // Clear unread
        this.markGroupAsRead(groupID)
      } catch (err) {
        console.error('Failed to fetch group messages:', err)
      }
    },
    async sendGroupMessage(groupID, body, file = null) {
      let imageURL = null
      if (file) {
        try {
          imageURL = await this.uploadImage(file)
        } catch (err) {
          console.error('Image upload failed, message not sent', err)
          throw err
        }
      }

      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
        throw new Error('COMM_LINK_OFFLINE')
      }

      const payload = {
        type: 'group_message',
        data: {
          group_id: groupID,
          body: body,
          image_url: imageURL
        }
      }
      this.socket.send(JSON.stringify(payload))
    },
    getUnreadCount(id, type = 'u') {
      return this.unreadCounts[`${type}_${id}`] || 0
    }
  },

  getters: {
    totalUnread: (state) => {
      return Object.values(state.unreadCounts).reduce((a, b) => a + b, 0)
    }
  }
})
