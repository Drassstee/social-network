<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const auth = useAuthStore()
const profile = ref(null)
const loading = ref(true)

const fetchProfile = async (id) => {
  loading.value = true
  try {
    const response = await fetch(`/api/v1/users/${id}`)
    if (!response.ok) throw new Error('Failed to fetch profile')
    profile.value = await response.json()
  } catch (e) { 
    console.error(e)
    profile.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchProfile(route.params.id))
watch(() => route.params.id, (newId) => fetchProfile(newId))

const togglePrivacy = async () => {
  if (!profile.value) return
  const newType = profile.value.user.profile_type === 'public' ? 'private' : 'public'
  try {
    const response = await fetch('/api/v1/users', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_type: newType })
    })
    if (response.ok) {
      profile.value.user.profile_type = newType
    }
  } catch (e) { console.error(e) }
}
</script>

<template>
  <div class="profile-view">
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="profile" class="profile-content">
      <div class="card-traditional profile-header">
        <div class="profile-cover">
          <div class="bg-seigaiha"></div>
        </div>
        
        <div class="profile-info-section">
          <div class="avatar-large">{{ profile.user.first_name[0] }}</div>
          <div class="user-details">
            <h1 class="fullname">{{ profile.user.first_name }} {{ profile.user.last_name }}</h1>
            <p class="nickname">@{{ profile.user.nickname || 'user' + profile.user.id }}</p>
          </div>
          
          <div class="profile-actions">
            <button 
              v-if="auth.user?.id == route.params.id" 
              @click="togglePrivacy" 
              class="btn-traditional ghost"
            >
              Set to {{ profile.user.profile_type === 'public' ? 'Private' : 'Public' }}
            </button>
            <button v-else class="btn-traditional">Follow</button>
          </div>
        </div>
        
        <div class="stats-row">
          <div class="stat-item">
            <span class="stat-value">{{ profile.followers?.length || 0 }}</span>
            <span class="stat-label">Followers</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ profile.following?.length || 0 }}</span>
            <span class="stat-label">Following</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ profile.posts?.length || 0 }}</span>
            <span class="stat-label">Posts</span>
          </div>
        </div>

        <div class="about-section">
          <h3>About Me</h3>
          <p>{{ profile.user.about_me || 'No description provided.' }}</p>
        </div>
      </div>

      <div class="profile-activity">
        <h2>Activity</h2>
        <div class="posts-list">
          <div v-if="profile.posts?.length === 0" class="no-posts">
            No posts to display.
          </div>
          <!-- Post cards loop here -->
        </div>
      </div>
    </div>
    <div v-else class="error-msg">
      <h2>Profile not found or could not be loaded.</h2>
      <p>Please double check the URL or try again later.</p>
    </div>
  </div>
</template>

<style scoped>
.profile-header {
  padding: 0;
  overflow: hidden;
  margin-bottom: 40px;
}

.profile-cover {
  height: 150px;
  background: var(--color-charcoal);
  position: relative;
  overflow: hidden;
}

.profile-info-section {
  display: flex;
  padding: 0 40px;
  margin-top: -60px;
  align-items: flex-end;
  gap: 30px;
  margin-bottom: 30px;
}

.avatar-large {
  width: 150px;
  height: 150px;
  background: var(--color-charcoal);
  color: var(--color-gold);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 5rem;
  font-family: 'Noto Serif JP', serif;
  border: 6px solid var(--color-washi-white);
  box-shadow: var(--shadow-japanese);
  z-index: 10;
}

.user-details {
  flex: 1;
  padding-bottom: 10px;
}

.fullname {
  font-size: 2.2rem;
  color: var(--color-charcoal);
}

.nickname {
  color: var(--color-vermilion);
  font-family: 'Noto Serif JP', serif;
}

.profile-actions {
  padding-bottom: 10px;
}

.btn-traditional.ghost {
  background: none;
  color: var(--color-vermilion);
  border: 2px solid var(--color-vermilion);
}

.stats-row {
  display: flex;
  justify-content: center;
  gap: 60px;
  padding: 20px 40px;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-charcoal);
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
}

.about-section {
  padding: 30px 40px;
}

.about-section h3 {
  margin-bottom: 10px;
  font-size: 1.3rem;
}

.profile-activity {
  padding: 0 10px;
}

.profile-activity h2 {
  margin-bottom: 25px;
}

.no-posts {
  text-align: center;
  padding: 40px;
  color: #888;
  font-style: italic;
}
</style>
