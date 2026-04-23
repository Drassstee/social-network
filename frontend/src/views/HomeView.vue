<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const posts = ref([])
const newPost = ref({
  body: '',
  privacy: 'public',
  allowed_users: []
})

// Mock fetching posts for now
onMounted(async () => {
  try {
    const response = await fetch('/api/v1/posts')
    if (response.ok) {
        const data = await response.json()
        posts.value = data.posts || []
    }
  } catch (e) { console.error('Failed to fetch posts:', e) }
})

const handleCreatePost = async () => {
  try {
    const resp = await fetch('/api/v1/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        author_id: String(auth.user?.id),
        content: newPost.value.body
      })
    })
    if (resp.ok) {
      const data = await resp.json()
      if (data.post) {
        // Since backend doesn't embed the full author object, inject the current user's profile
        // Handling both common JSON mapping cases for robustness
        data.post.author = {
          first_name: auth.user?.first_name || auth.user?.FirstName || 'Anonymous',
          last_name: auth.user?.last_name || auth.user?.LastName || ''
        }
        posts.value.unshift(data.post)
      }
    }
  } catch (e) { console.error('Failed to create post:', e) }
  finally { newPost.value.body = '' }
}

</script>

<template>
  <div class="home-view">
    <header class="view-header">
      <h1 class="view-title">Feed</h1>
      <p class="view-subtitle">Feed</p>
    </header>

    <div class="card-traditional create-post-card">
      <div class="user-info">
        <div class="avatar-placeholder">{{ auth.user?.first_name?.[0] || '?' }}</div>
        <span class="username">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</span>
      </div>
      
      <form @submit.prevent="handleCreatePost" class="post-form">
        <textarea 
          v-model="newPost.body" 
          placeholder="What is on your mind?" 
          class="input-traditional textarea"
          rows="3"
          required
        ></textarea>
        
        <div class="post-actions">
          <div class="privacy-select">
            <select v-model="newPost.privacy" class="input-traditional select-mini">
              <option value="public">🌐 Public</option>
              <option value="almost_private">👥 Followers</option>
              <option value="private">🔒 Private</option>
            </select>
          </div>
          <button type="submit" class="btn-traditional">Post</button>
        </div>
      </form>
    </div>

    <div class="posts-list">
      <div v-for="post in posts" :key="post.id" class="card-traditional post-card">
        <div class="post-header">
          <div class="avatar-placeholder">{{ post.author?.first_name?.[0] || 'U' }}</div>
          <div class="post-meta">
            <h3 class="author-name">{{ post.author?.first_name || post.author?.FirstName || 'User ' + (post.author_id || '') }} {{ post.author?.last_name || post.author?.LastName || '' }}</h3>
            <span class="post-date text-muted">{{ new Date(post.created_at).toLocaleDateString() }}</span>
          </div>
          <div class="privacy-badge">{{ post.privacy || 'public' }}</div>
        </div>
        
        <div class="post-body">
          <p>{{ post.content || post.body }}</p>
        </div>
        
        <div class="post-footer">
          <button class="action-btn">💬 Comment</button>
        </div>
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
  font-size: 2.5rem;
}

.view-subtitle {
  color: var(--color-vermilion);
  font-size: 1.2rem;
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

.avatar-placeholder {
  width: 45px;
  height: 45px;
  background: var(--color-charcoal);
  color: var(--color-paper);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Noto Serif JP', serif;
  font-weight: 700;
  border: 2px solid var(--color-gold);
}

.username {
  font-weight: 600;
}

.textarea {
  resize: none;
  font-size: 1.1rem;
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
  font-size: 0.9rem;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.post-meta {
  flex: 1;
}

.author-name {
  font-size: 1.1rem;
}

.text-muted {
  font-size: 0.8rem;
  color: #666;
}

.privacy-badge {
  font-size: 0.75rem;
  background: var(--color-paper);
  padding: 2px 8px;
  border-radius: 10px;
  border: 1px solid #ddd;
  text-transform: capitalize;
}

.post-body {
  margin-bottom: 20px;
  font-size: 1.1rem;
  color: var(--color-charcoal);
  line-height: 1.8;
}

.post-footer {
  border-top: 1px solid #eee;
  padding-top: 15px;
}

.action-btn {
  background: none;
  border: none;
  color: var(--color-charcoal);
  cursor: pointer;
  font-weight: 500;
  transition: color 0.3s;
}

.action-btn:hover {
  color: var(--color-vermilion);
}
</style>
