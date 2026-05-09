<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import UserAvatar from './UserAvatar.vue'

const props = defineProps({
  post: Object
})

const auth = useAuthStore()
const showComments = ref(false)
const newCommentContent = ref('')
const commentImage = ref(null)

const onFileChange = (e) => {
  commentImage.value = e.target.files[0]
}

const toggleComments = () => {
  showComments.value = !showComments.value
}

const handleCreateComment = async () => {
  if (!newCommentContent.value && !commentImage.value) return

  try {
    const formData = new FormData()
    formData.append('content', newCommentContent.value)
    formData.append('post_id', props.post.id)
    if (commentImage.value) {
      formData.append('image', commentImage.value)
    }

    const resp = await fetch('/api/v1/posts/comments', {
      method: 'POST',
      body: formData
    })
    
    if (resp.ok) {
      const data = await resp.json()
      if (!props.post.comments) props.post.comments = []
      
      // Add optimistic author data if missing from response
      if (data && !data.author) {
        data.author = {
          id: auth.user.id,
          first_name: auth.user.first_name,
          last_name: auth.user.last_name,
          avatar_url: auth.user.avatar_url
        }
      }
      
      props.post.comments.push(data)
      newCommentContent.value = ''
      commentImage.value = null
    }
  } catch (e) {
    console.error('Failed to create comment:', e)
  }
}
</script>

<template>
  <div class="card-retro post-card">
    <div class="post-header">
      <router-link :to="`/profile/${post.author_id}`" class="avatar-link">
        <UserAvatar 
          :url="post.author?.avatar_url" 
          :name="post.author?.first_name" 
          size="small" 
        />
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
      <button @click="toggleComments" class="action-btn">
        💬 Comment ({{ post.comments?.length || 0 }})
      </button>
    </div>

    <div v-if="showComments" class="comments-section">
      <div class="comments-list">
        <div v-for="comment in post.comments" :key="comment.id" class="comment-item">
          <UserAvatar 
            :url="comment.author?.avatar_url" 
            :name="comment.author?.first_name" 
            size="xsmall" 
          />
          <div class="comment-content">
            <span class="comment-author">{{ comment.author?.first_name }} {{ comment.author?.last_name }}</span>
            <p>{{ comment.content }}</p>
            <img v-if="comment.image_url" :src="comment.image_url" class="comment-image">
          </div>
        </div>
      </div>
      
      <div class="comment-input-area">
        <textarea 
          v-model="newCommentContent" 
          placeholder="Write a comment..." 
          class="input-retro mini-textarea"
        ></textarea>
        <div class="comment-actions">
          <label class="comment-file-label">
            🖼️
            <input type="file" @change="onFileChange" accept="image/*" class="hidden-input">
          </label>
          <button @click="handleCreateComment" class="btn-retro mini-btn">Reply</button>
        </div>
        <span v-if="commentImage" class="file-name-mini">{{ commentImage.name }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-card {
  margin-bottom: 20px;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.avatar-link {
  text-decoration: none;
}

.post-meta {
  flex-grow: 1;
}

.author-link {
  text-decoration: none;
}

.author-name {
  margin: 0;
  font-size: 1.1rem;
  color: var(--color-neon-yellow);
}

.post-date {
  font-size: 0.8rem;
}

.privacy-badge {
  font-size: 0.7rem;
  padding: 2px 6px;
  border: 1px solid var(--color-neon-magenta);
  color: var(--color-neon-magenta);
  text-transform: uppercase;
  font-family: 'VT323', monospace;
}

.post-body {
  margin-bottom: 15px;
}

.post-body p {
  font-size: 1.2rem;
  line-height: 1.4;
  margin-bottom: 10px;
  white-space: pre-wrap;
}

.post-image {
  width: 100%;
  border-radius: 4px;
  border: 1px solid var(--color-grid-line);
}

.post-footer {
  border-top: 1px solid var(--color-grid-line);
  padding-top: 10px;
}

.action-btn {
  background: none;
  border: none;
  color: var(--color-neon-cyan);
  cursor: pointer;
  font-family: 'VT323', monospace;
  font-size: 1.1rem;
  transition: all 0.2s;
}

.action-btn:hover {
  text-shadow: var(--shadow-neon-cyan);
}

.comments-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px dashed var(--color-grid-line);
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 15px;
}

.comment-item {
  display: flex;
  gap: 10px;
}

.comment-content {
  background: rgba(255, 255, 255, 0.05);
  padding: 8px 12px;
  border-radius: 8px;
  flex-grow: 1;
}

.comment-author {
  font-weight: bold;
  font-size: 0.9rem;
  color: var(--color-neon-magenta);
  display: block;
  margin-bottom: 4px;
}

.comment-content p {
  margin: 0;
  font-size: 1rem;
}

.comment-image {
  max-width: 200px;
  margin-top: 8px;
  border-radius: 4px;
}

.comment-input-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mini-textarea {
  min-height: 60px;
  font-size: 1rem;
}

.comment-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}

.comment-file-label {
  cursor: pointer;
  font-size: 1.2rem;
}

.hidden-input {
  display: none;
}

.mini-btn {
  padding: 4px 12px;
  font-size: 0.9rem;
}

.file-name-mini {
  font-size: 0.8rem;
  color: var(--color-neon-cyan);
  text-align: right;
}
</style>
