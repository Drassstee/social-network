<script setup>
import { onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useProfileStore } from '../stores/profile'
import { useUIStore } from '../stores/ui'
import UserAvatar from '../components/UserAvatar.vue'
import PostCard from '../components/PostCard.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'

const route = useRoute()
const auth = useAuthStore()
const profileStore = useProfileStore()
const ui = useUIStore()

const isMe = computed(() => {
  return auth.user && profileStore.profile?.user?.id == auth.user.id
})

const isPrivateAndNotFollowed = computed(() => {
  if (isMe.value) return false
  return profileStore.profile?.privacy === 'private' && profileStore.followStatus !== 'following'
})

const fetchProfile = async (id) => {
  try {
    await profileStore.fetchProfile(id, auth.user?.id)
  } catch (e) { 
    console.error(e)
  }
}

onMounted(() => fetchProfile(route.params.id))
watch(() => route.params.id, (newId) => fetchProfile(newId))

const handleFollow = async () => {
  try {
    const status = await profileStore.follow(route.params.id)
    ui.showToast(status === 'following' ? 'Following!' : 'Follow request sent!', 'success')
    fetchProfile(route.params.id)
  } catch (e) { console.error(e) }
}

const handleUnfollow = async () => {
  if (!window.confirm('DISCONNECT FROM THIS NODE? (UNFOLLOW)')) return
  try {
    await profileStore.unfollow(route.params.id)
    ui.showToast('Unfollowed', 'success')
    fetchProfile(route.params.id)
  } catch (e) { console.error(e) }
}

const togglePrivacy = async () => {
  if (!profileStore.profile) return
  const current = profileStore.profile.user
  const newType = profileStore.profile.privacy === 'public' ? 'private' : 'public'
  
  if (!window.confirm(`CHANGE PROTOCOL TO ${newType.toUpperCase()}?`)) return

  try {
    await profileStore.updateProfile({ 
      id: profileStore.profile.user.id,
      profile_type: newType
    })
    ui.showToast(`Profile is now ${newType}`, 'success')
    fetchProfile(route.params.id)
  } catch (e) { console.error(e) }
}
</script>

<template>
  <div class="profile-view">
    <div v-if="profileStore.loading && !profileStore.profile" class="profile-skeleton">
      <SkeletonLoader type="avatar" />
      <div class="skel-meta">
        <SkeletonLoader type="text" :count="2" />
      </div>
      <div class="skel-content">
        <SkeletonLoader type="card" :count="2" />
      </div>
    </div>

    <div v-else-if="profileStore.profile" class="profile-content">
      <div class="card-retro profile-header">
        <div class="profile-cover">
          <div class="cover-gradient"></div>
        </div>
        
        <div class="profile-info-section">
          <div class="avatar-container-large">
            <UserAvatar 
              :url="profileStore.profile.user.avatar_url" 
              :name="profileStore.profile.user.first_name" 
              size="large" 
            />
          </div>
          <div class="user-details">
            <h1 class="fullname">{{ profileStore.profile.user.first_name }} {{ profileStore.profile.user.last_name }}</h1>
            <p class="nickname">ID_{{ profileStore.profile.user.id }} // @{{ profileStore.profile.user.nickname || 'user' }}</p>
          </div>
          
          <div class="profile-actions">
            <button 
              v-if="isMe" 
              @click="togglePrivacy" 
              class="btn-retro ghost"
            >
              MODE: {{ profileStore.profile.privacy }}
            </button>
            <template v-else>
              <button v-if="profileStore.followStatus === 'none'" @click="handleFollow" class="btn-retro">CONNECT</button>
              <button v-else-if="profileStore.followStatus === 'pending'" class="btn-retro disabled" disabled>PENDING...</button>
              <button v-else @click="handleUnfollow" class="btn-retro ghost">DISCONNECT</button>
            </template>
          </div>
        </div>
        
        <div class="stats-row">
          <router-link :to="`/profile/${route.params.id}/followers`" class="stat-item linkable">
            <span class="stat-value">{{ profileStore.profile.followers?.length || 0 }}</span>
            <span class="stat-label">NODES_IN</span>
          </router-link>
          <router-link :to="`/profile/${route.params.id}/following`" class="stat-item linkable">
            <span class="stat-value">{{ profileStore.profile.following?.length || 0 }}</span>
            <span class="stat-label">NODES_OUT</span>
          </router-link>
          <div class="stat-item">
            <span class="stat-value">{{ profileStore.profile.posts?.length || 0 }}</span>
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
              <span class="info-value">{{ profileStore.profile.email }}</span>
            </div>
            <div v-if="profileStore.profile.dob" class="info-item">
              <span class="info-label">EXP_DATE:</span>
              <span class="info-value">{{ new Date(profileStore.profile.dob).toLocaleDateString() }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">PROTOCOL:</span>
              <span class="info-value text-capitalize">{{ profileStore.profile.privacy }}</span>
            </div>
          </div>
          <div class="bio-section">
            <h3>LOG_FILE</h3>
            <p>{{ profileStore.profile.about_me || 'NO DATA RECOVERED.' }}</p>
          </div>
        </div>
      </div>

      <div v-if="!isPrivateAndNotFollowed" class="profile-activity">
        <h2>TRANSMISSIONS</h2>
        <div class="posts-list">
          <div v-if="profileStore.profile.posts?.length === 0" class="no-posts">
            NO DATA RECORDED.
          </div>
          <PostCard v-else v-for="post in profileStore.profile.posts" :key="post.id" :post="post" />
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

.avatar-container-large {
  flex-shrink: 0;
  margin-top: -60px;
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
.profile-skeleton {
  display: flex;
  flex-direction: column;
  gap: 30px;
  align-items: center;
}

.skel-meta {
  width: 200px;
}

.skel-content {
  width: 100%;
  max-width: 600px;
}
</style>
