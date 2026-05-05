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
    router.push('/login')
  }
})
</script>

<template>
  <div class="layout-wrapper">
    <div class="bg-seigaiha"></div>
    
    <nav class="sidebar">
      <div class="logo">
        <h1 class="logo-text">Kizuna</h1>
        <div class="logo-accent">C</div>
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
            <span class="nav-role">Member</span>
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
  background: var(--color-charcoal);
  color: var(--color-paper);
  height: 100vh;
  position: sticky;
  top: 0;
  display: flex;
  flex-direction: column;
  padding: 40px 20px;
  border-right: 4px solid var(--color-gold);
}

.logo {
  text-align: center;
  margin-bottom: 60px;
}

.logo-text {
  color: var(--color-gold);
  font-size: 2.5rem;
  letter-spacing: 2px;
}

.logo-accent {
  font-size: 3rem;
  color: var(--color-vermilion);
  margin-top: -10px;
}

.nav-links {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.nav-item {
  color: var(--color-paper);
  text-decoration: none;
  font-family: 'Noto Serif JP', serif;
  font-size: 1.1rem;
  padding: 12px 20px;
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-item:hover, .router-link-active {
  background: rgba(212, 175, 55, 0.15);
  color: var(--color-gold);
  transform: translateX(5px);
}

.icon {
  font-size: 1.2rem;
}

.badge {
  background: var(--color-vermilion);
  color: white;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 10px;
  margin-left: auto;
  font-family: sans-serif;
  font-weight: bold;
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
  padding: 15px 20px;
  margin-top: auto;
  border-top: 1px solid rgba(212, 175, 55, 0.2);
}

.avatar-nav-img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid var(--color-gold);
  object-fit: cover;
}

.user-info-nav {
  display: flex;
  flex-direction: column;
}

.nav-username {
  font-weight: 700;
  font-size: 0.95rem;
  color: var(--color-gold);
}

.nav-role {
  font-size: 0.75rem;
  opacity: 0.7;
}

.main-content {
  flex: 1;
  padding: 40px;
  max-width: 1200px;
  margin: 0 auto;
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
    border-bottom: 4px solid var(--color-gold);
  }
  
  .logo {
    margin-bottom: 20px;
  }
}
</style>
