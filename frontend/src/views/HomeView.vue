<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { usePostStore } from '../stores/posts'
import { api } from '../api/api'
import UserAvatar from '../components/UserAvatar.vue'
import PostCard from '../components/PostCard.vue'

const auth = useAuthStore()
const postStore = usePostStore()
const newPost = ref({
  body: '',
  privacy: 'public',
  allowed_users: []
})
const postImage = ref(null)
const followers = ref([])
const feedbackMessage = ref('')
const feedbackType = ref('info')

const showFeedback = (msg, type = 'info') => {
  feedbackMessage.value = msg
  feedbackType.value = type
  setTimeout(() => { feedbackMessage.value = '' }, 3000)
}

const onFileChange = (e, type, postId = null) => {
  const file = e.target.files[0]
  if (type === 'post') postImage.value = file
}

const fetchFollowers = async () => {
  if (!auth.user?.id) return
  try {
    const data = await api.get(`/users/${auth.user.id}`)
    followers.value = data.followers || []
  } catch (e) { console.error('Failed to fetch followers:', e) }
}

onMounted(async () => {
  fetchFollowers()
  postStore.fetchPosts()
})

const handleCreatePost = async () => {
  try {
    await postStore.createPost(
      newPost.value.body,
      newPost.value.privacy,
      0,
      newPost.value.allowed_users,
      postImage.value
    )
    
    newPost.value.body = ''
    newPost.value.privacy = 'public'
    newPost.value.allowed_users = []
    postImage.value = null
    showFeedback('Post created!', 'success')
  } catch (e) { 
    console.error('Failed to create post:', e)
    showFeedback(e.message, 'error')
  }
}
</script>

<template>
  <div class="home-view">
    <div v-if="feedbackMessage" class="feedback-toast" :class="feedbackType">
      {{ feedbackMessage }}
    </div>
    <header class="view-header">
      <h1 class="view-title">Feed</h1>
      <p class="view-subtitle">Feed</p>
    </header>

    <div class="card-retro create-post-card">
      <div class="user-info">
        <UserAvatar :url="auth.user?.avatar_url" :name="auth.user?.first_name" size="small" />
        <span class="username">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</span>
      </div>
      
      <form @submit.prevent="handleCreatePost" class="post-form">
        <textarea 
          v-model="newPost.body" 
          placeholder="What is on your mind?" 
          class="input-retro textarea"
          rows="3"
        ></textarea>
        
        <div v-if="newPost.privacy === 'private'" class="allowed-users-selection">
          <label>Select followers to see this post:</label>
          <div class="followers-checklist">
            <label v-for="f in followers" :key="f.id" class="follower-checkbox">
              <input type="checkbox" :value="f.id" v-model="newPost.allowed_users">
              {{ f.first_name }} {{ f.last_name }}
            </label>
            <p v-if="followers.length === 0" class="text-muted italic">You have no followers to share with.</p>
          </div>
        </div>

        <div class="post-actions">
          <div class="action-left">
            <div class="privacy-select">
              <select v-model="newPost.privacy" class="input-retro select-mini">
                <option value="public">🌐 Public</option>
                <option value="almost_private">👥 Followers</option>
                <option value="private">🔒 Private</option>
              </select>
            </div>
            <label class="file-label">
              🖼️ Photo/GIF
              <input type="file" @change="e => onFileChange(e, 'post')" accept="image/*" class="hidden-input">
            </label>
            <span v-if="postImage" class="file-name">{{ postImage.name }}</span>
          </div>
          <button type="submit" class="btn-retro" :disabled="!newPost.body && !postImage">Post</button>
        </div>
      </form>
    </div>

    <div class="posts-list">
      <PostCard v-for="post in postStore.posts" :key="post.id" :post="post" />
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

.textarea {
  resize: none;
  font-size: 1.4rem;
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
  font-size: 1.3rem;
  color: var(--color-neon-cyan);
}

.avatar-link, .author-link {
  text-decoration: none;
  color: inherit;
  transition: all 0.2s;
}

.avatar-link:hover, .author-link:hover {
  filter: brightness(1.2);
}

.author-link:hover .author-name {
  color: var(--color-neon-magenta);
  text-shadow: 0 0 5px var(--color-neon-magenta);
}

.text-muted {
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
  padding: 10px;
  background: rgba(0,0,0,0.2);
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
  font-size: 1rem;
  color: var(--color-neon-magenta);
  display: block;
  margin-bottom: 2px;
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

.mini-btn {
  padding: 5px 15px;
  font-size: 0.8rem;
}

.italic { font-style: italic; }
.avatar-small-img {
  width: 45px;
  height: 45px;
  border: 2px solid var(--color-neon-magenta);
  object-fit: cover;
}

.avatar-xsmall-img {
  width: 28px;
  height: 28px;
  border: 1px solid var(--color-neon-magenta);
  object-fit: cover;
}

.file-label, .comment-file-label {
  cursor: pointer;
  background: rgba(255, 0, 255, 0.1);
  border: 1px solid var(--color-neon-magenta);
  padding: 6px 12px;
  font-size: 1rem;
  color: var(--color-neon-magenta);
  transition: all 0.3s;
}

.file-label:hover, .comment-file-label:hover {
  background: var(--color-neon-magenta);
  color: white;
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
</style>
