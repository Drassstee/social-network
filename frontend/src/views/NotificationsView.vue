<script setup>
import { onMounted } from 'vue'
import { useNotificationStore } from '../stores/notifications'

const store = useNotificationStore()

onMounted(async () => {
  await store.fetchNotifications()
  // Mark all as read when visiting the page? Or only when clicking specific ones?
  // Let's mark all as read for simplicity if the user stays on the page.
  setTimeout(() => store.markAllAsRead(), 2000)
})

const handleFollowResponse = async (followerId, status, notificationId) => {
  try {
    const response = await fetch('/api/v1/notifications/respond', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ follower_id: followerId, status: status })
    })
    if (response.ok) {
      if (notificationId) store.markAsRead(notificationId)
      store.fetchNotifications()
    }
  } catch (err) { console.error(err) }
}

const handleGroupInvitation = async (invitationId, accept, notificationId) => {
  try {
    const response = await fetch(`/api/v1/groups/invitations/${invitationId}/respond`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ accept: accept })
    })
    if (response.ok) {
      if (notificationId) store.markAsRead(notificationId)
      store.fetchNotifications()
    }
  } catch (err) { console.error(err) }
}

const handleJoinRequest = async (requestId, accept, notificationId) => {
  try {
    const response = await fetch(`/api/v1/groups/requests/${requestId}/respond`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ accept: accept })
    })
    if (response.ok) {
      if (notificationId) store.markAsRead(notificationId)
      store.fetchNotifications()
    }
  } catch (err) { console.error(err) }
}

const getIcon = (type) => {
  switch (type) {
    case 'follow_request': return '👤'
    case 'invite': return '📩'
    case 'request': return '🙋'
    case 'event': return '📅'
    case 'accept': return '✅'
    case 'decline': return '❌'
    default: return '🔔'
  }
}
</script>

<template>
  <div class="notifications-view">
    <header class="view-header">
      <h1 class="view-title">Notifications</h1>
      <div v-if="store.unreadCount > 0" class="unread-banner">
        You have {{ store.unreadCount }} new notifications
        <button @click="store.markAllAsRead" class="btn-text">Mark all as read</button>
      </div>
    </header>

    <div class="notifications-list">
      <div v-if="store.loading" class="loading">Loading notifications...</div>
      <div v-else-if="store.notifications.length === 0" class="card-traditional no-notifications">
        <p>No notifications yet.</p>
      </div>
      
      <div v-for="notif in store.notifications" :key="notif.id" 
           :class="['card-traditional notif-card', { 'unread': !notif.is_read }]">
        <div class="notif-content">
          <div class="notif-icon">{{ getIcon(notif.notification_type) }}</div>
          <div class="notif-text">
            <span v-if="notif.notification_type === 'follow_request'">
              <strong>{{ notif.actor_username }}</strong> sent you a follow request.
            </span>
            <span v-else-if="notif.notification_type === 'invite'">
              <strong>{{ notif.actor_username }}</strong> invited you to join a group.
            </span>
            <span v-else-if="notif.notification_type === 'request'">
              <strong>{{ notif.actor_username }}</strong> wants to join your group.
            </span>
            <span v-else-if="notif.notification_type === 'event'">
              New event in group: <strong>{{ notif.target_title || 'Group Event' }}</strong>
            </span>
            <span v-else-if="notif.notification_type === 'accept'">
              <strong>{{ notif.actor_username }}</strong> accepted your request.
            </span>
            <span v-else-if="notif.notification_type === 'decline'">
              <strong>{{ notif.actor_username }}</strong> declined your request.
            </span>
            <span v-else>
              {{ notif.actor_username }} triggered a {{ notif.notification_type }}.
            </span>
            <div class="notif-time">{{ new Date(notif.created_at).toLocaleString() }}</div>
          </div>
        </div>

        <div class="notif-actions" v-if="!notif.is_read || ['follow_request', 'invite', 'request'].includes(notif.notification_type)">
          <!-- Follow Request Actions -->
          <template v-if="notif.notification_type === 'follow_request'">
            <button @click="handleFollowResponse(notif.actor_id, 'accept', notif.id)" class="btn-traditional mini">Accept</button>
            <button @click="handleFollowResponse(notif.actor_id, 'decline', notif.id)" class="btn-traditional mini ghost">Decline</button>
          </template>

          <!-- Group Invite Actions -->
          <template v-else-if="notif.notification_type === 'invite'">
            <button @click="handleGroupInvitation(notif.target_id, true, notif.id)" class="btn-traditional mini">Join</button>
            <button @click="handleGroupInvitation(notif.target_id, false, notif.id)" class="btn-traditional mini ghost">Ignore</button>
          </template>

          <!-- Join Request Actions -->
          <template v-else-if="notif.notification_type === 'request'">
            <button @click="handleJoinRequest(notif.target_id, true, notif.id)" class="btn-traditional mini">Approve</button>
            <button @click="handleJoinRequest(notif.target_id, false, notif.id)" class="btn-traditional mini ghost">Deny</button>
          </template>

          <!-- Info Notifications -->
          <template v-else>
            <button v-if="!notif.is_read" @click="store.markAsRead(notif.id)" class="btn-text">Dismiss</button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-header {
  margin-bottom: 30px;
}

.unread-banner {
  background: rgba(212, 175, 55, 0.1);
  padding: 10px 20px;
  border-radius: var(--border-radius);
  margin-top: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
}

.btn-text {
  background: none;
  border: none;
  color: var(--color-gold);
  cursor: pointer;
  font-weight: 600;
  text-decoration: underline;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.no-notifications {
  text-align: center;
  padding: 40px;
  color: #888;
}

.notif-card {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  transition: all 0.3s ease;
  border-left: 4px solid transparent;
}

.notif-card.unread {
  background: rgba(212, 175, 55, 0.05);
  border-left-color: var(--color-gold);
}

.notif-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.notif-icon {
  font-size: 1.8rem;
  min-width: 40px;
  text-align: center;
}

.notif-text {
  font-size: 1rem;
}

.notif-time {
  font-size: 0.8rem;
  color: #888;
  margin-top: 4px;
}

.notif-actions {
  display: flex;
  gap: 10px;
}

.btn-traditional.mini {
  padding: 6px 15px;
  font-size: 0.85rem;
}

.btn-traditional.ghost {
  background: none;
  color: #666;
  border: 1px solid #ddd;
  box-shadow: none;
}

@media (max-width: 600px) {
  .notif-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  .notif-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
