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
      <h1 class="view-title">SYSTEM_ALERTS</h1>
      <div v-if="store.unreadCount > 0" class="unread-banner card-retro">
        <span>CRITICAL: {{ store.unreadCount }} NEW_ALERTS_DETECTED</span>
        <button @click="store.markAllAsRead" class="btn-retro mini">CLEAR_ALL</button>
      </div>
    </header>

    <div class="notifications-list">
      <div v-if="store.loading" class="loading">SCANNING_NETWORK...</div>
      <div v-else-if="store.notifications.length === 0" class="card-retro no-notifications">
        <p>STATUS_IDLE: NO_PENDING_NOTIFICATIONS.</p>
      </div>
      
      <div v-for="notif in store.notifications" :key="notif.id" 
           :class="['card-retro notif-card', { 'unread': !notif.is_read }]">
        <div class="notif-content">
          <div class="notif-icon">{{ getIcon(notif.notification_type) }}</div>
          <div class="notif-text">
            <span v-if="notif.message" class="notif-message-text">{{ notif.message }}</span>
            <template v-else>
              <span v-if="notif.notification_type === 'follow_request'">
                <strong class="actor">{{ notif.actor_username }}</strong> INITIATED_FOLLOW_SEQUENCE.
              </span>
              <span v-else-if="notif.notification_type === 'invite'">
                <strong class="actor">{{ notif.actor_username }}</strong> DISPATCHED_INVITE TO: <strong>{{ notif.target_title }}</strong>.
              </span>
              <span v-else-if="notif.notification_type === 'request'">
                <strong class="actor">{{ notif.actor_username }}</strong> REQUESTS ACCESS TO: <strong>{{ notif.target_title }}</strong>.
              </span>
              <span v-else-if="notif.notification_type === 'event'">
                NEW_MISSION_OBJECTIVE IN <strong>{{ notif.target_title }}</strong>.
              </span>
              <span v-else-if="notif.notification_type === 'accept'">
                <strong class="actor">{{ notif.actor_username }}</strong> JOINED: <strong>{{ notif.target_title }}</strong>.
              </span>
              <span v-else-if="notif.notification_type === 'decline'">
                <strong class="actor">{{ notif.actor_username }}</strong> DECLINED ACCESS TO: <strong>{{ notif.target_title }}</strong>.
              </span>
              <span v-else>
                {{ notif.actor_username }} TRIGGERED_{{ notif.notification_type.toUpperCase() }}.
              </span>
            </template>
            <div class="notif-time">TIMESTAMP: {{ new Date(notif.created_at).toLocaleString() }}</div>
          </div>

        </div>

        <div class="notif-actions" v-if="!notif.is_read || ['follow_request', 'invite', 'request'].includes(notif.notification_type)">
          <!-- Follow Request Actions -->
          <template v-if="notif.notification_type === 'follow_request'">
            <button @click="handleFollowResponse(notif.actor_id, 'accept', notif.id)" class="btn-retro mini">ALLOW</button>
            <button @click="handleFollowResponse(notif.actor_id, 'decline', notif.id)" class="btn-retro mini ghost">DENY</button>
          </template>

          <!-- Group Invite Actions -->
          <template v-else-if="notif.notification_type === 'invite'">
            <button @click="handleGroupInvitation(notif.target_id, true, notif.id)" class="btn-retro mini">ACCEPT</button>
            <button @click="handleGroupInvitation(notif.target_id, false, notif.id)" class="btn-retro mini ghost">IGNORE</button>
          </template>

          <!-- Join Request Actions -->
          <template v-else-if="notif.notification_type === 'request'">
            <button @click="handleJoinRequest(notif.target_id, true, notif.id)" class="btn-retro mini">APPROVE</button>
            <button @click="handleJoinRequest(notif.target_id, false, notif.id)" class="btn-retro mini ghost">DENY</button>
          </template>

          <!-- Info Notifications -->
          <template v-else>
            <button v-if="!notif.is_read" @click="store.markAsRead(notif.id)" class="btn-retro mini ghost">DISMISS</button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notifications-view {
  max-width: 900px;
  margin: 0 auto;
}

.view-header {
  margin-bottom: 40px;
}

.view-title {
  font-size: 2.5rem;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
  margin-bottom: 20px;
}

.unread-banner {
  background: rgba(255, 255, 0, 0.1);
  border-color: var(--color-neon-yellow);
  padding: 15px 25px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-family: 'VT323', monospace;
  font-size: 1.4rem;
  color: var(--color-neon-yellow);
  box-shadow: 5px 5px 0 var(--color-neon-yellow);
}

.loading {
  color: var(--color-neon-cyan);
  font-family: 'Press Start 2P', cursive;
  font-size: 0.8rem;
  text-align: center;
  padding: 40px;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  50% { opacity: 0.5; }
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-notifications {
  text-align: center;
  padding: 60px;
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}

.notif-card {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-left-width: 8px;
  transition: all 0.2s;
}

.notif-card.unread {
  background: rgba(0, 255, 255, 0.05);
  border-left-color: var(--color-neon-cyan);
}

.notif-content {
  display: flex;
  align-items: center;
  gap: 25px;
}

.notif-icon {
  font-size: 2.2rem;
  min-width: 50px;
  text-align: center;
  filter: drop-shadow(0 0 5px currentColor);
}

.notif-text {
  font-family: 'VT323', monospace;
  font-size: 1.3rem;
  color: white;
}

.actor {
  color: var(--color-neon-magenta);
  text-shadow: 0 0 5px var(--color-neon-magenta);
}

.notif-time {
  font-size: 0.9rem;
  color: var(--color-neon-yellow);
  margin-top: 5px;
}

.notif-actions {
  display: flex;
  gap: 15px;
}

.btn-retro.mini {
  padding: 8px 15px;
  font-size: 0.7rem;
}

.btn-retro.ghost {
  border-color: var(--color-grid-line);
  color: #888;
}

.btn-retro.ghost:hover {
  border-color: var(--color-neon-magenta);
  color: var(--color-neon-magenta);
}

@media (max-width: 768px) {
  .notif-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }
  .notif-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
