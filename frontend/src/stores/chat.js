import { defineStore } from 'pinia'

export const useChatStore = defineStore('chat', {
  state: () => ({
    socket: null,
    onlineUsers: [],
    messages: {}, // userID -> array of messages
    activeChatUser: null,
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
        console.log('Chat connected')
        this.fetchOnlineUsers()
      }

      this.socket.onmessage = (event) => {
        const msg = JSON.parse(event.data)
        this.handleSocketMessage(msg)
      }

      this.socket.onclose = () => {
        this.connected = false
        this.socket = null
        console.log('Chat disconnected')
        // Reconnect after a delay
        setTimeout(() => this.connect(), 3000)
      }
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
        case 'typing':
          // Handle typing indicator
          break
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
    addMessage(msg) {
      const otherID = msg.receiver_id === this.activeChatUser?.id ? msg.receiver_id : msg.sender_id
      if (!this.messages[otherID]) {
        this.messages[otherID] = []
      }
      this.messages[otherID].push(msg)
    },
    sendMessage(receiverID, body) {
      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return

      const payload = {
        type: 'private_message',
        data: {
          receiver_id: receiverID,
          body: body
        }
      }
      this.socket.send(JSON.stringify(payload))
    }
  }
})
