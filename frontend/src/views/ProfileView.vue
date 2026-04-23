<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const auth = useAuthStore()
const profile = ref(null)
const loading = ref(true)

const followStatus = ref('none') // 'none', 'pending', 'accept'
const isMe = ref(false)

const fetchProfile = async (id) => {
  loading.value = true
  isMe.value = auth.user?.id == id
  try {
    const response = await fetch(`/api/v1/users/${id}`)
    if (!response.ok) throw new Error('Failed to fetch profile')
    profile.value = await response.json()
    // In a real app we'd fetch follow status specifically, 
    // but here we can check if I'm in followers list if it's public
    if (profile.value.followers?.some(f => f.id === auth.user?.id)) {
      followStatus.value = 'accept'
    }
  } catch (e) { 
    console.error(e)
    profile.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchProfile(route.params.id))
watch(() => route.params.id, (newId) => fetchProfile(newId))

const handleFollow = async () => {
  if (!profile.value) return
  try {
    const response = await fetch('/api/v1/follow', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ following_id: Number(route.params.id) })
    })
    if (response.ok) {
      const data = await response.json()
      followStatus.value = data.status || 'pending'
      alert(data.status === 'accept' ? 'Following!' : 'Follow request sent!')
      fetchProfile(route.params.id)
    }
  } catch (e) { console.error(e) }
}

const handleUnfollow = async () => {
  if (!profile.value) return
  try {
    const response = await fetch('/api/v1/unfollow', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ following_id: Number(route.params.id) })
    })
    if (response.ok) {
      followStatus.value = 'none'
      alert('Unfollowed')
      fetchProfile(route.params.id)
    }
  } catch (e) { console.error(e) }
}

const togglePrivacy = async () => {
  if (!profile.value) return
  const newType = profile.value.user.profile_type === 'public' ? 'private' : 'public'
  try {
    const response = await fetch('/api/v1/users', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profile_type: newType, email: profile.value.user.email, first_name: profile.value.user.first_name, last_name: profile.value.user.last_name, dob: profile.value.user.dob })
    })
    if (response.ok) {
      profile.value.user.profile_type = newType
      alert(`Profile is now ${newType}`)
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
          <div class="avatar-placeholder avatar-large">{{ profile.user.first_name[0] }}</div>
          <div class="user-details">
            <h1 class="fullname">{{ profile.user.first_name }} {{ profile.user.last_name }}</h1>
            <p class="nickname">@{{ profile.user.nickname || 'user' + profile.user.id }}</p>
          </div>
          
          <div class="profile-actions">
            <button 
              v-if="isMe" 
              @click="togglePrivacy" 
              class="btn-traditional ghost"
            >
              Set to {{ profile.user.profile_type === 'public' ? 'Private' : 'Public' }}
            </button>
            <template v-else>
              <button v-if="followStatus === 'none'" @click="handleFollow" class="btn-traditional">Follow</button>
              <button v-else-if="followStatus === 'pending'" class="btn-traditional disabled" disabled>Request Pending</button>
              <button v-else @click="handleUnfollow" class="btn-traditional ghost">Unfollow</button>
            </template>
          </div>
        </div>
        
        <div class="stats-row">
          <router-link :to="`/profile/${route.params.id}/followers`" class="stat-item linkable">
            <span class="stat-value">{{ profile.followers?.length || 0 }}</span>
            <span class="stat-label">Followers</span>
          </router-link>
          <router-link :to="`/profile/${route.params.id}/following`" class="stat-item linkable">
            <span class="stat-value">{{ profile.following?.length || 0 }}</span>
            <span class="stat-label">Following</span>
          </router-link>
          <div class="stat-item">
            <span class="stat-value">{{ profile.posts?.length || 0 }}</span>
            <span class="stat-label">Posts</span>
          </div>
        </div>

        <div class="about-section">
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">Email:</span>
              <span class="info-value">{{ profile.user.email }}</span>
            </div>
            <div v-if="profile.user.dob" class="info-item">
              <span class="info-label">Birthday:</span>
              <span class="info-value">{{ new Date(profile.user.dob).toLocaleDateString() }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Privacy:</span>
              <span class="info-value text-capitalize">{{ profile.user.profile_type }}</span>
            </div>
          </div>
          <div class="bio-section">
            <h3>About Me</h3>
            <p>{{ profile.user.about_me || 'No description provided.' }}</p>
          </div>
        </div>
      </div>

      <div class="profile-activity">
        <h2>Activity</h2>
        <div class="posts-list">
          <div v-if="profile.posts?.length === 0" class="no-posts">
            No posts to display.
          </div>
          <div v-else v-for="post in profile.posts" :key="post.id" class="card-traditional post-card">
            <div class="post-header">
              <div class="avatar-placeholder small">{{ profile.user.first_name[0] }}</div>
              <div class="post-meta">
                <h3 class="author-name">{{ profile.user.first_name }} {{ profile.user.last_name }}</h3>
                <span class="post-date text-muted">{{ new Date(post.created_at).toLocaleDateString() }}</span>
              </div>
            </div>
            <div class="post-body">
              <p>{{ post.content }}</p>
            </div>
          </div>
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
  font-size: 5rem;
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

.stat-item.linkable {
  cursor: pointer;
  text-decoration: none;
  transition: transform 0.2s;
}

.stat-item.linkable:hover {
  transform: translateY(-3px);
}

.stat-item.linkable:hover .stat-label {
  color: var(--color-vermilion);
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

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
  background: rgba(212, 175, 55, 0.05);
  padding: 20px;
  border-radius: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 0.8rem;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.info-value {
  font-weight: 600;
  color: var(--color-charcoal);
}

.text-capitalize {
  text-transform: capitalize;
}

.bio-section h3 {
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
