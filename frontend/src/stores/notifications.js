import { defineStore } from 'pinia'

export const useNotificationStore = defineStore('notifications', {
  state: () => ({
    notifications: [],
    unreadCount: 0,
    loading: false
  }),
  actions: {
    async fetchNotifications() {
      this.loading = true
      try {
        const response = await fetch('/api/v1/notifications/list')
        const data = await response.json()
        this.notifications = data || []
      } catch (err) {
        console.error('Failed to fetch notifications:', err)
      } finally {
        this.loading = false
      }
    },
    async fetchUnreadCount() {
      try {
        const response = await fetch('/api/v1/notifications/unread-count')
        const data = await response.json()
        this.unreadCount = data.count || 0
      } catch (err) {
        console.error('Failed to fetch unread count:', err)
      }
    },
    async markAsRead(id) {
      try {
        const response = await fetch(`/api/v1/notifications/${id}/read`, { method: 'POST' })
        if (response.ok) {
          const notif = this.notifications.find(n => n.id === id)
          if (notif && !notif.is_read) {
            notif.is_read = true
            this.unreadCount = Math.max(0, this.unreadCount - 1)
          }
        }
      } catch (err) {
        console.error('Failed to mark notification as read:', err)
      }
    },
    async markAllAsRead() {
      try {
        const response = await fetch('/api/v1/notifications/read-all', { method: 'POST' })
        if (response.ok) {
          this.notifications.forEach(n => n.is_read = true)
          this.unreadCount = 0
        }
      } catch (err) {
        console.error('Failed to mark all as read:', err)
      }
    },
    addNotification(notification) {
      this.notifications.unshift(notification)
      if (!notification.is_read) {
        this.unreadCount++
      }
    },
    incrementUnread() {
      this.unreadCount++
    }
  }
})
