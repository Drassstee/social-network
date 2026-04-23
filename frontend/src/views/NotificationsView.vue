<script setup>
import { onMounted } from 'vue'
import { useNotificationStore } from '../stores/notifications'

const store = useNotificationStore()

onMounted(() => {
  store.fetchNotifications()
})

const getNotificationText = (notif) => {
  switch (notif.type) {
    case 'follow_request': return 'sent you a follow request.'
    case 'group_invitation': return 'invited you to join'
    case 'join_request': return 'requested to join your group'
    case 'event_created': return 'created a new event in'
    default: return 'sent you a notification.'
  }
}
</script>

<template>
  <div class="notifications-view">
    <header class="view-header">
      <h1 class="view-title">Notifications</h1>
      <p class="view-subtitle">Notifications</p>
    </header>

    <div class="notifications-list">
      <div v-if="store.loading" class="loading">Loading notifications...</div>
      <div v-else-if="store.notifications.length === 0" class="card-traditional no-notifications">
        <p>No notifications yet.</p>
      </div>
      <div v-for="notif in store.notifications" :key="notif.id" class="card-traditional notif-card">
        <div class="notif-content">
          <div class="notif-icon">🔔</div>
          <div class="notif-text">
            <strong>User {{ notif.sender_id }}</strong> {{ getNotificationText(notif) }}
            <span v-if="notif.group_id" class="group-name"> Group {{ notif.group_id }}</span>
          </div>
        </div>
        <div class="notif-actions" v-if="notif.type === 'follow_request' || notif.type === 'group_invitation'">
          <button class="btn-traditional mini">Accept</button>
          <button class="btn-traditional mini ghost">Decline</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-notifications {
  text-align: center;
  padding: 40px;
  color: #888;
}

.notif-card {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.notif-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.notif-icon {
  font-size: 1.5rem;
}

.group-name {
  color: var(--color-vermilion);
  font-weight: 600;
}

.notif-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.btn-traditional.mini {
  padding: 5px 15px;
  font-size: 0.85rem;
}

.btn-traditional.ghost {
  background: none;
  color: #666;
  border: 1px solid #ddd;
  box-shadow: none;
}
</style>
