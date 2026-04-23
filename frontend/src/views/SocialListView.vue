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
      <h1 class="view-title">{{ title }}</h1>
      <p class="view-subtitle">{{ title }} List</p>
    </header>

    <div class="card-traditional">
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="users.length === 0" class="no-data">
        No users found in this list.
      </div>
      <div v-else class="users-grid">
        <router-link v-for="user in users" :key="user.id" :to="`/profile/${user.id}`" class="user-card-mini">
          <div class="avatar-placeholder">{{ user.first_name[0] }}</div>
          <div class="user-info-mini">
            <span class="fullname">{{ user.first_name }} {{ user.last_name }}</span>
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
  font-size: 2.2rem;
}

.view-subtitle {
  color: var(--color-vermilion);
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.user-card-mini {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s ease;
}

.user-card-mini:hover {
  border-color: var(--color-gold);
  background: rgba(212, 175, 55, 0.05);
  transform: translateY(-2px);
}

.fullname {
  font-weight: 600;
  font-size: 1.1rem;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #888;
  font-style: italic;
}
</style>
