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
const postImage = ref(null)
const followers = ref([])
const showComments = ref({})
const newCommentContent = ref({})
const commentImages = ref({})
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
  else commentImages.value[postId] = file
}

const fetchFollowers = async () => {
  if (!auth.user?.id) return
  try {
    const response = await fetch(`/api/v1/users/${auth.user.id}`)
    if (response.ok) {
      const data = await response.json()
      followers.value = data.followers || []
    }
  } catch (e) { console.error('Failed to fetch followers:', e) }
}

onMounted(async () => {

  fetchFollowers()
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
    const formData = new FormData()
    formData.append('content', newPost.value.body)
    formData.append('privacy', newPost.value.privacy)
    if (newPost.value.privacy === 'private') {
      formData.append('allowed_users', JSON.stringify(newPost.value.allowed_users))
    }
    if (postImage.value) {
      formData.append('image', postImage.value)
    }

    const resp = await fetch('/api/v1/posts', {
      method: 'POST',
      body: formData
    })
    
    if (resp.ok) {
      const data = await resp.json()
      if (data.post) {
        data.post.author = {
          first_name: auth.user?.first_name || 'Me',
          last_name: auth.user?.last_name || ''
        }
        posts.value.unshift(data.post)
        newPost.value.body = ''
        newPost.value.privacy = 'public'
        newPost.value.allowed_users = []
        postImage.value = null
        showFeedback('Post created!', 'success')
      }
    }
  } catch (e) { console.error('Failed to create post:', e) }
}

const toggleComments = (postId) => {
  showComments.value[postId] = !showComments.value[postId]
}

const handleCreateComment = async (postId) => {
  const content = newCommentContent.value[postId]
  const image = commentImages.value[postId]
  if (!content && !image) return

  try {
    const formData = new FormData()
    formData.append('post_id', postId)
    formData.append('content', content || '')
    if (image) formData.append('image', image)

    const resp = await fetch('/api/v1/comments', {
      method: 'POST',
      body: formData
    })

    if (resp.ok) {
      const comment = await resp.json()
      const post = posts.value.find(p => p.id === postId)
      if (post) {
        if (!post.comments) post.comments = []
        post.comments.push(comment)
      }
      newCommentContent.value[postId] = ''
      commentImages.value[postId] = null
    }
  } catch (e) { console.error('Failed to create comment:', e) }
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
              <select v-model="newPost.privacy" class="input-traditional select-mini">
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
          <button type="submit" class="btn-traditional" :disabled="!newPost.body && !postImage">Post</button>
        </div>
      </form>
    </div>

    <div class="posts-list">
      <div v-for="post in posts" :key="post.id" class="card-traditional post-card">
        <div class="post-header">
          <router-link :to="`/profile/${post.author_id}`" class="avatar-link">
            <img v-if="post.author?.avatar_url" :src="post.author.avatar_url" class="avatar-small-img" />
            <div v-else class="avatar-placeholder">{{ post.author?.first_name?.[0] || 'U' }}</div>
          </router-link>
          <div class="post-meta">
            <router-link :to="`/profile/${post.author_id}`" class="author-link">
              <h3 class="author-name">{{ post.author?.first_name || 'User' }} {{ post.author?.last_name || '' }}</h3>
            </router-link>
            <span class="post-date text-muted">{{ new Date(post.created_at).toLocaleDateString() }}</span>
          </div>
          <div class="privacy-badge">{{ post.privacy || 'public' }}</div>
        </div>
        
        <div class="post-body">
          <p v-if="post.content">{{ post.content }}</p>
          <img v-if="post.image_url" :src="post.image_url" class="post-image" alt="Post content">
        </div>
        
        <div class="post-footer">
          <button @click="toggleComments(post.id)" class="action-btn">💬 Comment ({{ post.comments?.length || 0 }})</button>
        </div>

        <div v-if="showComments[post.id]" class="comments-section">
          <div class="comments-list">
            <div v-for="comment in post.comments" :key="comment.id" class="comment-item">
               <img v-if="comment.author?.avatar_url" :src="comment.author.avatar_url" class="avatar-xsmall-img" />
               <div v-else class="avatar-placeholder xsmall">{{ comment.author?.first_name?.[0] || 'U' }}</div>
              <div class="comment-content">
                <span class="comment-author">{{ comment.author?.first_name }} {{ comment.author?.last_name }}</span>
                <p>{{ comment.content }}</p>
                <img v-if="comment.image_url" :src="comment.image_url" class="comment-image">
              </div>
            </div>
          </div>
          
          <div class="comment-input-area">
            <textarea 
              v-model="newCommentContent[post.id]" 
              placeholder="Write a comment..." 
              class="input-traditional mini-textarea"
            ></textarea>
            <div class="comment-actions">
              <label class="comment-file-label">
                🖼️
                <input type="file" @change="e => onFileChange(e, 'comment', post.id)" accept="image/*" class="hidden-input">
              </label>
              <button @click="handleCreateComment(post.id)" class="btn-traditional mini-btn">Reply</button>
            </div>
            <span v-if="commentImages[post.id]" class="file-name-mini">{{ commentImages[post.id].name }}</span>
          </div>
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

.avatar-link, .author-link {
  text-decoration: none;
  color: inherit;
  transition: opacity 0.2s;
}

.avatar-link:hover, .author-link:hover {
  opacity: 0.8;
}

.author-link:hover .author-name {
  color: var(--color-vermilion);
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

.action-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.file-label, .comment-file-label {
  cursor: pointer;
  background: #eee;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: background 0.3s;
}

.file-label:hover {
  background: var(--color-gold);
  color: white;
}

.hidden-input {
  display: none;
}

.file-name, .file-name-mini {
  font-size: 0.8rem;
  color: #666;
}

.allowed-users-selection {
  margin-top: 15px;
  padding: 15px;
  background: #f9f9f9;
  border-radius: 8px;
  border: 1px dashed #ccc;
}

.followers-checklist {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}

.follower-checkbox {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.9rem;
  background: white;
  padding: 4px 10px;
  border-radius: 20px;
  border: 1px solid #ddd;
  cursor: pointer;
}

.post-image {
  max-width: 100%;
  border-radius: 8px;
  margin-top: 15px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.comments-section {
  padding-top: 20px;
  margin-top: 15px;
  border-top: 1px solid #eee;
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
  background: #f1f1f1;
  padding: 8px 15px;
  border-radius: 18px;
  flex: 1;
}

.comment-author {
  font-weight: 700;
  font-size: 0.9rem;
  display: block;
  margin-bottom: 2px;
}

.comment-image {
  max-width: 200px;
  border-radius: 4px;
  margin-top: 8px;
}

.comment-input-area {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mini-textarea {
  min-height: 40px;
  font-size: 0.95rem;
  padding: 10px;
  border-radius: 12px;
}

.comment-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}

.mini-btn {
  padding: 5px 15px;
  font-size: 0.9rem;
}

.italic { font-style: italic; }
.avatar-small-img {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  border: 2px solid var(--color-gold);
  object-fit: cover;
}

.avatar-xsmall-img {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid var(--color-gold);
  object-fit: cover;
}
</style>
