import { defineStore } from 'pinia'

export const useNotificationStore = defineStore('notifications', {
  state: () => ({
    notifications: [],
    loading: false
  }),
  actions: {
    async fetchNotifications() {
      this.loading = true
      try {
        const response = await fetch('/api/v1/notifications')
        const data = await response.json()
        this.notifications = data || []
      } catch (err) {
        console.error('Failed to fetch notifications:', err)
      } finally {
        this.loading = false
      }
    },
    addNotification(notification) {
      this.notifications.unshift(notification)
    }
  }
})
