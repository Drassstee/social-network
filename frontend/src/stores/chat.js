import { defineStore } from 'pinia'
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
        this.fetchOnlineUsers()
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
    async fetchOnlineUsers() {
      try {
        const response = await fetch('/api/v1/chat/online')
        const data = await response.json()
        this.onlineUsers = data
      } catch (err) {
        console.error('Failed to fetch online users:', err)
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
      const index = this.onlineUsers.findIndex(u => u.id === data.user_id)
      if (data.online && index === -1) {
        this.onlineUsers.push({ id: data.user_id, username: data.username })
      } else if (!data.online && index !== -1) {
        this.onlineUsers.splice(index, 1)
      }
    },
    async addMessage(msg) {
      // Lazy import auth to get current user ID
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      const myID = authStore.user?.id
      
      const otherID = msg.sender_id === myID ? msg.receiver_id : msg.sender_id
      
      if (!this.messages[otherID]) {
        this.messages[otherID] = []
      }
      this.messages[otherID].push(msg)

      // Increment unread if not active chat
      if (msg.sender_id !== myID && (!this.activeChatUser || this.activeChatUser.id !== otherID)) {
        const key = `u_${otherID}`
        this.unreadCounts[key] = (this.unreadCounts[key] || 0) + 1
      }
    },
    sendMessage(receiverID, body, imageURL = null) {
      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return

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
        const response = await fetch(`/api/v1/chat/messages?user_id=${userID}&limit=50&offset=0`)
        const data = await response.json()
        this.messages[userID] = data || []
        
        // Clear unread
        delete this.unreadCounts[`u_${userID}`]
      } catch (err) {
        console.error('Failed to fetch messages:', err)
      }
    },
    async addGroupMessage(msg) {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      const myID = authStore.user?.id
      
      const groupID = msg.group_id
      if (!this.groupMessages[groupID]) {
        this.groupMessages[groupID] = []
      }
      
      if (!this.groupMessages[groupID].some(m => m.id === msg.id)) {
        this.groupMessages[groupID].push(msg)
      }

      // Increment unread if not active group
      if (Number(msg.sender_id) !== Number(myID) && Number(this.activeGroupID) !== Number(groupID)) {
        const key = `g_${groupID}`
        this.unreadCounts[key] = (this.unreadCounts[key] || 0) + 1
      }
      return msg
    },
    async fetchGroupMessages(groupID) {
      try {
        const response = await fetch(`/api/v1/groups/${groupID}/messages`)
        const data = await response.json()
        this.groupMessages[groupID] = data || []
        
        // Clear unread
        delete this.unreadCounts[`g_${groupID}`]
      } catch (err) {
        console.error('Failed to fetch group messages:', err)
      }
    },
    sendGroupMessage(groupID, body, imageURL = null) {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return

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
