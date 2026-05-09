<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import UserAvatar from '../components/UserAvatar.vue'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const auth = useAuthStore()
const profile = ref(null)
const loading = ref(true)

const followStatus = ref('none')
const isMe = computed(() => {
  return auth.user && profile.value?.user?.id == auth.user.id
})
const feedbackMessage = ref('')
const feedbackType = ref('info')

const showFeedback = (msg, type = 'info') => {
  feedbackMessage.value = msg
  feedbackType.value = type
  setTimeout(() => { feedbackMessage.value = '' }, 3000)
}

const isPrivateAndNotFollowed = computed(() => {
  if (isMe.value) return false
  return profile.value?.user?.profile_type === 'private' && followStatus.value !== 'accept'
})

const fetchProfile = async (id) => {
  loading.value = true
  try {
    const response = await fetch(`/api/v1/users/${id}`)
    if (!response.ok) throw new Error('Failed to fetch profile')
    const data = await response.json()
    profile.value = data
    
    // Determine follow status from backend response
    if (data.followers?.some(f => f.id === auth.user?.id)) {
      followStatus.value = 'accept'
    } else {
      followStatus.value = 'none'
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
      showFeedback(data.status === 'accept' ? 'Following!' : 'Follow request sent!', 'success')
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
      showFeedback('Unfollowed', 'success')
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
      body: JSON.stringify({ 
        profile_type: newType, 
        email: profile.value.user.email, 
        first_name: profile.value.user.first_name, 
        last_name: profile.value.user.last_name, 
        dob: profile.value.user.dob 
      })
    })
    if (response.ok) {
      profile.value.user.profile_type = newType
      showFeedback(`Profile is now ${newType}`, 'success')
    }
  } catch (e) { console.error(e) }
}

</script>


<template>
  <div class="profile-view">
    <div v-if="feedbackMessage" class="feedback-toast" :class="feedbackType">
      {{ feedbackMessage }}
    </div>
    <div v-if="loading" class="loading">SCANNING DATABASE...</div>
    <div v-else-if="profile" class="profile-content">
      <div class="card-retro profile-header">
        <div class="profile-cover">
          <div class="cover-gradient"></div>
        </div>
        
        <div class="profile-info-section">
          <div class="avatar-container-large">
            <UserAvatar 
              :url="profile.user.avatar_url" 
              :name="profile.user.first_name" 
              size="large" 
            />
          </div>
          <div class="user-details">
            <h1 class="fullname">{{ profile.user.first_name }} {{ profile.user.last_name }}</h1>
            <p class="nickname">ID_{{ profile.user.id }} // @{{ profile.user.nickname || 'user' }}</p>
          </div>
          
          <div class="profile-actions">
            <button 
              v-if="isMe" 
              @click="togglePrivacy" 
              class="btn-retro ghost"
            >
              MODE: {{ profile.user.profile_type }}
            </button>
            <template v-else>
              <button v-if="followStatus === 'none'" @click="handleFollow" class="btn-retro">CONNECT</button>
              <button v-else-if="followStatus === 'pending'" class="btn-retro disabled" disabled>PENDING...</button>
              <button v-else @click="handleUnfollow" class="btn-retro ghost">DISCONNECT</button>
            </template>
          </div>
        </div>
        
        <div class="stats-row">
          <router-link :to="`/profile/${route.params.id}/followers`" class="stat-item linkable">
            <span class="stat-value">{{ profile.followers?.length || 0 }}</span>
            <span class="stat-label">NODES_IN</span>
          </router-link>
          <router-link :to="`/profile/${route.params.id}/following`" class="stat-item linkable">
            <span class="stat-value">{{ profile.following?.length || 0 }}</span>
            <span class="stat-label">NODES_OUT</span>
          </router-link>
          <div class="stat-item">
            <span class="stat-value">{{ profile.posts?.length || 0 }}</span>
            <span class="stat-label">DATA_PACKS</span>
          </div>
        </div>

        <!-- Privacy Gating -->
        <div v-if="isPrivateAndNotFollowed" class="private-overlay">
          <div class="private-icon">🚫</div>
          <h2>ENCRYPTED ACCESS</h2>
          <p>ESTABLISH CONNECTION TO DECRYPT DATA.</p>
        </div>

        <div v-else class="about-section">
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">NET_ADDR:</span>
              <span class="info-value">{{ profile.user.email }}</span>
            </div>
            <div v-if="profile.user.dob" class="info-item">
              <span class="info-label">EXP_DATE:</span>
              <span class="info-value">{{ new Date(profile.user.dob).toLocaleDateString() }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">PROTOCOL:</span>
              <span class="info-value text-capitalize">{{ profile.user.profile_type }}</span>
            </div>
          </div>
          <div class="bio-section">
            <h3>LOG_FILE</h3>
            <p>{{ profile.user.about_me || 'NO DATA RECOVERED.' }}</p>
          </div>
        </div>
      </div>

      <div v-if="!isPrivateAndNotFollowed" class="profile-activity">
        <h2>TRANSMISSIONS</h2>
        <div class="posts-list">
          <div v-if="profile.posts?.length === 0" class="no-posts">
            NO DATA RECORDED.
          </div>
          <PostCard v-else v-for="post in profile.posts" :key="post.id" :post="post" />
        </div>
      </div>
    </div>
    <div v-else class="error-msg">
      <h2>ERROR 404: NODE_NOT_FOUND</h2>
      <p>TARGET DISCONNECTED OR MOVED TO RESTRICTED SECTOR.</p>
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
  height: 180px;
  background: var(--color-deep-purple);
  position: relative;
  overflow: hidden;
  border-bottom: 2px solid var(--color-neon-magenta);
}

.cover-gradient {
  width: 100%;
  height: 100%;
  background: linear-gradient(0deg, var(--color-neon-magenta) 0%, transparent 100%);
  opacity: 0.3;
}

.profile-info-section {
  display: flex;
  padding: 0 40px;
  position: relative;
  z-index: 2;
  align-items: flex-end;
  gap: 30px;
  margin-bottom: 30px;
}

.avatar-large {
  width: 150px;
  height: 150px;
  font-size: 5rem;
}

.user-details {
  flex: 1;
  padding-bottom: 10px;
}

.fullname {
  font-size: 2.2rem;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
  margin-bottom: 5px;
}

.nickname {
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
  font-size: 1.4rem;
}

.profile-actions {
  padding-bottom: 10px;
}

.btn-retro.ghost {
  background: rgba(255, 0, 255, 0.1);
}

.stats-row {
  display: flex;
  justify-content: center;
  gap: 40px;
  padding: 20px 40px;
  border-top: 1px solid var(--color-grid-line);
  border-bottom: 1px solid var(--color-grid-line);
  background: rgba(0, 255, 255, 0.05);
}

.stat-item {
  text-align: center;
}

.stat-item.linkable {
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s;
}

.stat-item.linkable:hover {
  transform: scale(1.1);
  text-shadow: 0 0 10px var(--color-neon-cyan);
}

.stat-value {
  display: block;
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--color-neon-yellow);
}

.stat-label {
  font-size: 0.9rem;
  color: var(--color-neon-cyan);
}

.about-section {
  padding: 30px 40px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
  background: rgba(0, 255, 255, 0.05);
  padding: 20px;
  border: 1px solid var(--color-neon-cyan);
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 0.9rem;
  color: var(--color-neon-magenta);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.info-value {
  font-weight: 600;
  color: var(--color-neon-yellow);
  font-size: 1.2rem;
}

.text-capitalize {
  text-transform: capitalize;
}

.bio-section h3 {
  margin-bottom: 10px;
  font-size: 1.5rem;
  color: var(--color-neon-cyan);
}

.bio-section p {
  font-size: 1.3rem;
  color: white;
  line-height: 1.4;
}

.profile-activity {
  padding: 0 10px;
}

.profile-activity h2 {
  font-size: 2rem;
  margin-bottom: 20px;
  text-align: center;
}

.post-card {
  margin-bottom: 25px;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.post-meta {
  flex: 1;
}

.author-name {
  font-size: 1.3rem;
  color: var(--color-neon-cyan);
}

.post-date {
  font-size: 1rem;
  color: rgba(0, 255, 255, 0.6);
}

.privacy-badge {
  font-size: 0.9rem;
  background: rgba(255, 0, 255, 0.1);
  color: var(--color-neon-magenta);
  padding: 2px 10px;
  border: 1px solid var(--color-neon-magenta);
  text-transform: uppercase;
}

.post-body {
  margin-bottom: 20px;
  font-size: 1.3rem;
  color: var(--color-neon-yellow);
  line-height: 1.4;
  padding: 15px;
  background: rgba(0,0,0,0.3);
}

.post-image {
  max-width: 100%;
  max-height: 400px;
  object-fit: contain;
  border: 1px solid var(--color-neon-cyan);
  margin-top: 15px;
  display: block;
}

.post-footer {
  border-top: 1px solid var(--color-grid-line);
  padding-top: 15px;
}

.action-btn {
  background: none;
  border: none;
  color: var(--color-neon-cyan);
  cursor: pointer;
  font-weight: 700;
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  transition: all 0.2s;
}

.action-btn:hover {
  color: var(--color-neon-magenta);
  text-shadow: 0 0 5px var(--color-neon-magenta);
}

.comments-section {
  padding-top: 20px;
  margin-top: 15px;
  border-top: 1px solid var(--color-grid-line);
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 20px;
}

.comment-item {
  display: flex;
  gap: 10px;
}

.comment-content {
  background: rgba(0, 255, 255, 0.05);
  border: 1px solid var(--color-neon-cyan);
  padding: 8px 15px;
  flex: 1;
}

.comment-author {
  font-weight: 700;
  font-size: 1.1rem;
  color: var(--color-neon-magenta);
  display: block;
}

.comment-image {
  max-width: 200px;
  border: 1px solid var(--color-neon-cyan);
  margin-top: 8px;
}

.comment-input-area {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mini-textarea {
  min-height: 40px;
  font-size: 1.2rem;
  padding: 10px;
}

.comment-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}

.comment-file-label {
  cursor: pointer;
  background: rgba(255, 0, 255, 0.1);
  border: 1px solid var(--color-neon-magenta);
  padding: 6px 12px;
  color: var(--color-neon-magenta);
  transition: all 0.3s;
}

.comment-file-label:hover {
  background: var(--color-neon-magenta);
  color: white;
}

.hidden-input {
  display: none;
}

.file-name-mini {
  font-size: 0.9rem;
  color: rgba(0, 255, 255, 0.6);
}

.mini-btn {
  padding: 5px 15px;
  font-size: 0.9rem;
}

.private-overlay {
  padding: 60px 40px;
  text-align: center;
  background: rgba(0, 0, 0, 0.4);
  border-top: 2px solid var(--color-neon-magenta);
}

.private-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  text-shadow: 0 0 20px var(--color-neon-magenta);
}

.private-overlay h2 {
  font-size: 2rem;
  margin-bottom: 10px;
  color: var(--color-neon-magenta);
}

.private-overlay p {
  color: var(--color-neon-cyan);
  font-size: 1.4rem;
}

.italic { font-style: italic; }
.avatar-container-large {
  flex-shrink: 0;
  margin-top: -60px;
}

.avatar-large-img {
  width: 150px;
  height: 150px;
  border: 4px solid var(--color-neon-magenta);
  box-shadow: 0 0 20px var(--color-neon-magenta);
  object-fit: cover;
}

.avatar-small-img {
  width: 45px;
  height: 45px;
  border: 2px solid var(--color-neon-magenta);
  object-fit: cover;
}

.loading, .error-msg {
  text-align: center;
  padding: 100px 20px;
  font-family: 'Press Start 2P', cursive;
  color: var(--color-neon-magenta);
  text-shadow: var(--shadow-neon-magenta);
}

.error-msg h2 {
  font-size: 2rem;
  margin-bottom: 20px;
}

.error-msg p {
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
  color: var(--color-neon-cyan);
}
</style>
