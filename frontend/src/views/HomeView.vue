<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { usePostStore } from '../stores/posts'
import { useProfileStore } from '../stores/profile'
import { useUIStore } from '../stores/ui'
import PostCard from '../components/PostCard.vue'
import CreatePost from '../components/CreatePost.vue'
import SkeletonLoader from '../components/SkeletonLoader.vue'

const auth = useAuthStore()
const postStore = usePostStore()
const profileStore = useProfileStore()
const ui = useUIStore()

onMounted(async () => {
  if (auth.user?.id) {
    profileStore.fetchMyFollowers(auth.user.id)
  }
  postStore.fetchPosts()
})
</script>

<template>
  <div class="home-view">
    <header class="view-header">
      <h1 class="view-title">TRANSMISSIONS</h1>
      <p class="view-subtitle">DECRYPTING_GLOBAL_FEED...</p>
    </header>

    <CreatePost />

    <div class="posts-list">
      <SkeletonLoader v-if="postStore.loading && postStore.posts.length === 0" :count="3" />
      
      <template v-else>
        <PostCard v-for="post in postStore.posts" :key="post.id" :post="post" />
        <div v-if="postStore.posts.length > 0" class="end-of-feed">
          --- END_OF_TRANSMISSION ---
        </div>
      </template>
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
  font-size: 1.5rem;
  font-family: 'VT323', monospace;
}

.create-post-card {
  margin-bottom: 40px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.username {
  font-weight: 700;
  color: var(--color-neon-yellow);
  font-size: 1.2rem;
}

.post-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 15px;
}

.select-mini {
  padding: 5px 10px;
  width: auto;
  font-size: 1.1rem;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.allowed-users-selection {
  margin-top: 15px;
  padding: 15px;
  background: rgba(0, 255, 255, 0.05);
  border: 1px dashed var(--color-neon-cyan);
}

.follower-checkbox {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 1rem;
  background: var(--color-dark-bg);
  padding: 4px 10px;
  border: 1px solid var(--color-neon-magenta);
  cursor: pointer;
  color: var(--color-neon-cyan);
}
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.end-of-feed {
  text-align: center;
  padding: 40px;
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  letter-spacing: 2px;
  opacity: 0.6;
}
</style>
