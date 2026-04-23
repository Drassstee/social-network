<script setup>
import { onMounted } from 'vue'
import { useNotificationStore } from '../stores/notifications'

const store = useNotificationStore()

onMounted(() => {
  store.fetchNotifications()
})

const handleResponse = async (followerId, status) => {
  try {
    const response = await fetch('/api/v1/notifications/respond', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ follower_id: followerId, status: status })
    })
    if (response.ok) {
      alert(`Request ${status}ed`)
      store.fetchNotifications()
    }
  } catch (err) { console.error(err) }
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
      <div v-for="user in store.notifications" :key="user.id" class="card-traditional notif-card">
        <div class="notif-content">
          <div class="notif-icon">👤</div>
          <div class="notif-text">
            <strong>{{ user.first_name }} {{ user.last_name }}</strong> sent you a follow request.
          </div>
        </div>
        <div class="notif-actions">
          <button @click="handleResponse(user.id, 'accept')" class="btn-traditional mini">Accept</button>
          <button @click="handleResponse(user.id, 'decline')" class="btn-traditional mini ghost">Decline</button>
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
