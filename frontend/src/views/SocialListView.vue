<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const users = ref([])
const loading = ref(true)
const title = ref('')

const fetchList = async () => {
  loading.value = true
  const type = route.path.includes('followers') ? 'followers' : 'following'
  title.value = type.charAt(0) + type.slice(1)
  
  try {
    const response = await fetch(`/api/v1/users/${route.params.id}`)
    if (response.ok) {
      const data = await response.json()
      users.value = data[type] || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchList)
watch(() => route.path, fetchList)
</script>

<template>
  <div class="social-list-view">
    <header class="view-header">
      <h1 class="view-title">{{ title.toUpperCase() }}_DATABASE</h1>
      <p class="view-subtitle">SCANNING_USER_NETWORK...</p>
    </header>

    <div class="card-retro">
      <div v-if="loading" class="loading">FETCHING_DATA...</div>
      <div v-else-if="users.length === 0" class="no-data">
        STATUS: NO_USERS_DETECTED_IN_THIS_NODE.
      </div>
      <div v-else class="users-grid">
        <router-link v-for="user in users" :key="user.id" :to="`/profile/${user.id}`" class="user-card-mini card-retro">
          <div class="avatar-sm">{{ user.first_name[0] }}</div>
          <div class="user-info-mini">
            <span class="fullname">{{ user.first_name }} {{ user.last_name }}</span>
            <span class="uid">UID_{{ user.id.toString().padStart(4, '0') }}</span>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-header {
  margin-bottom: 40px;
  text-align: center;
}

.view-title {
  font-size: 3rem;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
}

.view-subtitle {
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}

.loading {
  text-align: center;
  padding: 40px;
  color: var(--color-neon-cyan);
  font-family: 'Press Start 2P', cursive;
  font-size: 0.8rem;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  padding: 20px;
}

.user-card-mini {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  text-decoration: none;
  color: white;
  transition: all 0.2s;
  background: rgba(31, 11, 53, 0.5);
}

.user-card-mini:hover {
  transform: scale(1.05);
  box-shadow: 0 0 15px var(--color-neon-cyan);
  border-color: var(--color-neon-cyan);
}

.avatar-sm {
  width: 45px;
  height: 45px;
  background: var(--color-dark-bg);
  border: 1px solid var(--color-neon-magenta);
  color: var(--color-neon-magenta);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Press Start 2P', cursive;
  font-size: 1.2rem;
  box-shadow: 2px 2px 0 var(--color-neon-cyan);
}

.user-info-mini {
  display: flex;
  flex-direction: column;
}

.fullname {
  font-weight: 700;
  font-size: 1.2rem;
  font-family: 'VT323', monospace;
  color: var(--color-neon-cyan);
}

.uid {
  font-size: 0.8rem;
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
}

.no-data {
  text-align: center;
  padding: 60px;
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}
</style>
