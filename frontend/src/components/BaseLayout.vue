<script setup>
import { useAuthStore } from '../stores/auth'
import { useChatStore } from '../stores/chat'
import { useNotificationStore } from '../stores/notifications'
import { useRouter } from 'vue-router'
import { onMounted, watch } from 'vue'

const auth = useAuthStore()
const chat = useChatStore()
const notifStore = useNotificationStore()
const router = useRouter()

const handleLogout = async () => {
  await auth.logout()
  router.push('/login')
}

onMounted(async () => {
  await auth.checkSession()
  if (auth.isAuthenticated) {
    chat.connect()
    notifStore.fetchUnreadCount()
  }
})

watch(() => auth.isAuthenticated, (val) => {
  if (val) {
    chat.connect()
    notifStore.fetchUnreadCount()
  } else {
    chat.disconnect()
    notifStore.unreadCount = 0
    notifStore.notifications = []
    router.push('/login')
  }
})
</script>

<template>
  <div class="layout-wrapper">
    <div class="bg-synthwave-grid"></div>
    <div class="crt-overlay"></div>
    <nav class="sidebar">
      <div class="logo">
        <h1 class="logo-text">Network</h1>
      </div>
      
      <div v-if="auth.isAuthenticated" class="nav-links">
        <router-link to="/" class="nav-item">
          <span class="icon">🏠</span> Feed
        </router-link>
        <router-link :to="`/profile/${auth.user?.id}`" class="nav-item">
          <span class="icon">👤</span> Profile
        </router-link>
        <router-link to="/groups" class="nav-item">
          <span class="icon">👥</span> Groups
        </router-link>
        <router-link to="/chat" class="nav-item">
          <span class="icon">💬</span> Chat
          <span v-if="chat.totalUnread > 0" class="badge chat-badge">{{ chat.totalUnread }}</span>
        </router-link>
        <router-link to="/notifications" class="nav-item">
          <span class="icon">🔔</span> Notifications
          <span v-if="notifStore.unreadCount > 0" class="badge">{{ notifStore.unreadCount }}</span>
        </router-link>
        
        <div class="user-profile-summary">
          <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" class="avatar-nav-img" />
          <div v-else class="avatar-placeholder small">{{ auth.user?.first_name?.[0] || 'U' }}</div>
          <div class="user-info-nav">
            <span class="nav-username">{{ auth.user?.first_name }}</span>
            <span class="nav-role">OPERATOR</span>
          </div>
        </div>

        <button @click="handleLogout" class="nav-item logout-btn">
          <span class="icon">🚪</span> Logout
        </button>
      </div>
      
      <div v-else class="nav-links">
        <router-link to="/login" class="nav-item">Login</router-link>
        <router-link to="/register" class="nav-item">Register</router-link>
      </div>
    </nav>
    
    <main class="main-content">
      <div class="content-container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<style scoped>
.layout-wrapper {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 280px;
  background: rgba(11, 12, 16, 0.95);
  color: var(--color-neon-cyan);
  height: 100vh;
  position: sticky;
  top: 0;
  display: flex;
  flex-direction: column;
  padding: 40px 20px;
  border-right: 4px solid var(--color-neon-magenta);
  box-shadow: 10px 0 20px rgba(255, 0, 255, 0.1);
  z-index: 10;
}

.logo {
  text-align: center;
  margin-bottom: 60px;
}

.logo-text {
  color: var(--color-neon-cyan);
  font-size: 2.2rem;
  letter-spacing: 2px;
  text-shadow: var(--shadow-neon-cyan);
}

.logo-accent {
  font-family: 'Press Start 2P', cursive;
  font-size: 2.5rem;
  color: var(--color-neon-magenta);
  text-shadow: var(--shadow-neon-magenta);
  margin-top: -5px;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.nav-item {
  color: var(--color-neon-cyan);
  text-decoration: none;
  font-family: 'VT323', monospace;
  font-size: 1.4rem;
  padding: 10px 15px;
  border: 1px solid transparent;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-item:hover, .router-link-active {
  background: rgba(0, 255, 255, 0.1);
  border-color: var(--color-neon-cyan);
  color: var(--color-neon-yellow);
  transform: translateX(10px);
  text-shadow: 0 0 5px var(--color-neon-yellow);
}

.icon {
  font-size: 1.4rem;
  filter: drop-shadow(0 0 5px var(--color-neon-magenta));
}

.badge {
  background: var(--color-neon-magenta);
  color: white;
  font-size: 0.8rem;
  padding: 2px 8px;
  border: 1px solid white;
  margin-left: auto;
  font-family: 'VT323', monospace;
}

.logout-btn {
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  margin-top: 20px;
}

.user-profile-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px;
  margin-top: auto;
  border: 1px solid var(--color-neon-cyan);
  background: rgba(0, 255, 255, 0.05);
}

.avatar-nav-img {
  width: 40px;
  height: 40px;
  border: 2px solid var(--color-neon-magenta);
  object-fit: cover;
}

.user-info-nav {
  display: flex;
  flex-direction: column;
}

.nav-username {
  font-weight: 700;
  font-size: 1.1rem;
  color: var(--color-neon-yellow);
}

.nav-role {
  font-size: 0.8rem;
  color: var(--color-neon-cyan);
  opacity: 0.8;
}

.main-content {
  flex: 1;
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
}

.content-container {
  max-width: 800px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .layout-wrapper {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    position: relative;
    padding: 20px;
    border-right: none;
    border-bottom: 4px solid var(--color-neon-magenta);
  }
  
  .logo {
    margin-bottom: 20px;
  }
}
</style>
