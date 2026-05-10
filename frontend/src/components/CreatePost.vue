<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { usePostStore } from '../stores/posts'
import { useProfileStore } from '../stores/profile'
import { useUIStore } from '../stores/ui'
import UserAvatar from './UserAvatar.vue'

const props = defineProps({
  groupId: {
    type: [Number, String],
    default: 0
  },
  initialPrivacy: {
    type: String,
    default: 'public'
  }
})

const auth = useAuthStore()
const postStore = usePostStore()
const profileStore = useProfileStore()
const ui = useUIStore()

const newPost = ref({
  body: '',
  privacy: props.initialPrivacy,
  allowed_users: []
})
const postImage = ref(null)
const imagePreview = ref(null)
const isSubmitting = ref(false)

const onFileChange = (e) => {
  const file = e.target.files[0]
  if (!file) return
  postImage.value = file
  
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target.result
  }
  reader.readAsDataURL(file)
}

const handleCreatePost = async () => {
  if (!newPost.value.body && !postImage.value) return
  
  isSubmitting.value = true
  try {
    await postStore.createPost(
      newPost.value.body,
      newPost.value.privacy,
      props.groupId,
      newPost.value.allowed_users,
      postImage.value
    )
    
    newPost.value.body = ''
    newPost.value.privacy = props.initialPrivacy
    newPost.value.allowed_users = []
    postImage.value = null
    ui.showToast('Post created!', 'success')
  } catch (e) { 
    console.error('Failed to create post:', e)
    ui.showToast(e.message, 'error')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="card-retro create-post-card">
    <div class="user-info">
      <UserAvatar :url="auth.user?.avatar_url" :name="auth.user?.first_name" size="small" />
      <span class="username">{{ auth.user?.first_name }} {{ auth.user?.last_name }}</span>
    </div>
    
    <form @submit.prevent="handleCreatePost" class="post-form">
      <textarea 
        v-model="newPost.body" 
        placeholder="What is on your mind?" 
        class="textarea-retro"
        rows="3"
      ></textarea>
      
      <div v-if="postImage" class="image-preview-container">
        <img :src="imagePreview" alt="Preview" class="image-preview">
        <button @click.prevent="postImage = null" class="btn-remove">×</button>
      </div>
      
      <div v-if="newPost.privacy === 'private'" class="allowed-users-selection">
        <label>Select followers to see this post:</label>
        <div class="followers-checklist">
          <label v-for="f in profileStore.myFollowers" :key="f.id" class="follower-checkbox">
            <input type="checkbox" :value="f.id" v-model="newPost.allowed_users">
            {{ f.first_name }} {{ f.last_name }}
          </label>
          <p v-if="profileStore.myFollowers.length === 0" class="text-muted italic">You have no followers to share with.</p>
        </div>
      </div>

      <div class="post-actions">
        <div class="action-left">
          <div class="privacy-select" v-if="!groupId">
            <select v-model="newPost.privacy" class="input-retro select-mini">
              <option value="public">🌐 Public</option>
              <option value="almost_private">👥 Followers</option>
              <option value="private">🔒 Private</option>
            </select>
          </div>
          <label class="file-label">
            🖼️ Photo/GIF
            <input type="file" @change="onFileChange" accept="image/*" class="hidden-input">
          </label>
          <span v-if="postImage" class="file-name">{{ postImage.name }}</span>
        </div>
        <button type="submit" class="btn-retro" :disabled="isSubmitting || (!newPost.body && !postImage)">
          {{ isSubmitting ? 'UPLOADING...' : 'Post' }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
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

.allowed-users-selection {
  margin-top: 15px;
  padding: 15px;
  background: rgba(0, 255, 255, 0.05);
  border: 1px dashed var(--color-neon-cyan);
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
  font-size: 1rem;
  background: var(--color-dark-bg);
  padding: 4px 10px;
  border: 1px solid var(--color-neon-magenta);
  cursor: pointer;
  color: var(--color-neon-cyan);
}

.file-name {
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
  font-size: 0.9rem;
  margin-left: 10px;
}

.text-muted {
  color: #aaa;
  font-size: 0.9rem;
}

.action-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.image-preview-container {
  position: relative;
  margin-top: 15px;
  max-width: 200px;
  border: 2px solid var(--color-neon-cyan);
  border-radius: 4px;
  overflow: hidden;
}

.image-preview {
  width: 100%;
  height: auto;
  display: block;
}

.btn-remove {
  position: absolute;
  top: 5px;
  right: 5px;
  background: var(--color-neon-magenta);
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  box-shadow: 0 0 10px rgba(255, 0, 255, 0.5);
}
</style>
