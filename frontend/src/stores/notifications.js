import { defineStore } from 'pinia'
import { api } from '../api/api'

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
        const data = await api.get('/notifications/list')
        this.notifications = data || []
      } catch (err) {
        console.error('Failed to fetch notifications:', err)
      } finally {
        this.loading = false
      }
    },
    async fetchUnreadCount() {
      try {
        const data = await api.get('/notifications/unread-count')
        this.unreadCount = data.count || 0
      } catch (err) {
        console.error('Failed to fetch unread count:', err)
      }
    },
    async markAsRead(id) {
      try {
        await api.post(`/notifications/${id}/read`)
        const notif = this.notifications.find(n => n.id === id)
        if (notif && !notif.is_read) {
          notif.is_read = true
          this.unreadCount = Math.max(0, this.unreadCount - 1)
        }
      } catch (err) {
        console.error('Failed to mark notification as read:', err)
      }
    },
    async markAllAsRead() {
      try {
        await api.post('/notifications/read-all')
        this.notifications.forEach(n => n.is_read = true)
        this.unreadCount = 0
      } catch (err) {
        console.error('Failed to mark all as read:', err)
      }
    },
    async respondToFollow(followerId, status) {
      return api.post('/notifications/respond', { follower_id: followerId, status: status })
    },
    async respondToGroupInvitation(invitationId, accept) {
      return api.post(`/groups/invitations/${invitationId}/respond`, { accept: accept })
    },
    async respondToJoinRequest(requestId, accept) {
      return api.post(`/groups/requests/${requestId}/respond`, { accept: accept })
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
