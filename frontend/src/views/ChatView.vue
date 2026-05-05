<script setup>
import { ref, computed, nextTick, watch, onMounted } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'

const chat = useChatStore()
const auth = useAuthStore()
const newMessage = ref('')
const showEmojiPicker = ref(false)
const messagesArea = ref(null)
const imageFile = ref(null)

const activeChatMessages = computed(() => {
  if (!chat.activeChatUser) return []
  return chat.messages[chat.activeChatUser.id] || []
})

const selectUser = (user) => {
  chat.activeChatUser = user
  chat.fetchMessages(user.id)
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesArea.value) {
      messagesArea.value.scrollTop = messagesArea.value.scrollHeight
    }
  })
}

watch(activeChatMessages, () => {
  scrollToBottom()
}, { deep: true })

const onSelectEmoji = (emoji) => {
  newMessage.value += emoji.i
}

const handleFileChange = (e) => {
  imageFile.value = e.target.files[0]
}


const handleSendMessage = async () => {
  if (!newMessage.value.trim() && !imageFile.value) return
  
  let imageURL = null
  if (imageFile.value) {
    const formData = new FormData()
    formData.append('image', imageFile.value)
    try {
      const resp = await fetch('/api/v1/chat/upload', {
        method: 'POST',
        body: formData
      })
      if (resp.ok) {
        const data = await resp.json()
        imageURL = data.url
      }
    } catch (err) { console.error('Upload failed:', err) }
  }

  chat.sendMessage(chat.activeChatUser.id, newMessage.value, imageURL)

  newMessage.value = ''
  imageFile.value = null
  showEmojiPicker.value = false
}

onMounted(() => {
  if (chat.activeChatUser) {
    chat.fetchMessages(chat.activeChatUser.id)
  }
})
</script>

<template>
  <div class="chat-container card-traditional">
    <div class="users-sidebar">
      <h3 class="sidebar-title">Online Members</h3>
      <div class="users-list">
        <div 
          v-for="user in chat.onlineUsers" 
          :key="user.id"
          class="user-item"
          :class="{ active: chat.activeChatUser?.id === user.id }"
          @click="selectUser(user)"
        >
          <div class="avatar-sm">{{ user.username?.[0] || '?' }}</div>
          <span class="username">{{ user.username }}</span>
          <div class="online-indicator"></div>
        </div>
      </div>
    </div>
    
    <div class="chat-main">
      <template v-if="chat.activeChatUser">
        <div class="chat-header">
          <h3>{{ chat.activeChatUser.username }}</h3>
        </div>
        
        <div class="messages-area" ref="messagesArea">
          <div 
            v-for="msg in activeChatMessages" 
            :key="msg.id"
            class="message-wrapper"
            :class="{ sent: msg.sender_id === auth.user.id }"
          >
            <div class="message-bubble">
              <img v-if="msg.image_url" :src="msg.image_url" class="chat-img" />
              <p v-if="msg.body">{{ msg.body }}</p>
              <span class="time">{{ new Date(msg.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</span>
            </div>
          </div>
        </div>
        
        <div class="input-controls">
          <div v-if="showEmojiPicker" class="emoji-picker-container">
            <EmojiPicker :native="true" @select="onSelectEmoji" />
          </div>
          
          <form @submit.prevent="handleSendMessage" class="message-input-form">
            <button type="button" @click="showEmojiPicker = !showEmojiPicker" class="icon-btn">😊</button>
            <label class="icon-btn file-label">
              🖼️
              <input type="file" @change="handleFileChange" accept="image/*" hidden />
            </label>
            <div class="input-wrapper">
              <input 
                v-model="newMessage" 
                type="text" 
                placeholder="Type a message..." 
                class="input-traditional"
                @focus="showEmojiPicker = false"
              />
              <div v-if="imageFile" class="file-preview">📎 {{ imageFile.name }}</div>
            </div>
            <button type="submit" class="btn-traditional">Send</button>
          </form>
        </div>
      </template>
      <div v-else class="no-chat-selected">
        <div class="no-chat-icon">💬</div>
        <p>Select a member to start chatting</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-container {
  display: flex;
  height: 75vh;
  padding: 0;
  overflow: hidden;
  border: 1px solid rgba(0,0,0,0.1);
}

.users-sidebar {
  width: 250px;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
  background: white;
}

.sidebar-title {
  padding: 20px;
  border-bottom: 1px solid #eee;
  font-size: 1.1rem;
  background: var(--color-charcoal);
  color: var(--color-gold);
}

.users-list {
  flex: 1;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 15px 20px;
  cursor: pointer;
  transition: background 0.3s;
}

.user-item:hover {
  background: #f9f9f9;
}

.user-item.active {
  background: rgba(188, 0, 45, 0.05);
  border-left: 4px solid var(--color-vermilion);
}

.avatar-sm {
  width: 35px;
  height: 35px;
  background: var(--color-charcoal);
  color: var(--color-paper);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}

.online-indicator {
  width: 10px;
  height: 10px;
  background: #4caf50;
  border-radius: 50%;
  margin-left: auto;
  box-shadow: 0 0 5px rgba(76, 175, 80, 0.5);
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--color-washi-white);
}

.chat-header {
  padding: 15px 25px;
  border-bottom: 1px solid #eee;
  background: white;
  box-shadow: 0 2px 5px rgba(0,0,0,0.02);
}

.messages-area {
  flex: 1;
  padding: 25px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 15px;
  scroll-behavior: smooth;
}

.message-wrapper {
  display: flex;
  flex-direction: column;
}

.message-wrapper.sent {
  align-items: flex-end;
}

.message-bubble {
  max-width: 70%;
  padding: 12px 18px;
  border-radius: 18px;
  background: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  position: relative;
}

.sent .message-bubble {
  background: var(--color-vermilion);
  color: white;
}

.chat-img {
  max-width: 100%;
  border-radius: 10px;
  margin-bottom: 8px;
}

.time {
  font-size: 0.65rem;
  opacity: 0.6;
  display: block;
  margin-top: 5px;
  text-align: right;
}

.input-controls {
  position: relative;
  border-top: 1px solid #eee;
  background: white;
}

.emoji-picker-container {
  position: absolute;
  bottom: 100%;
  right: 20px;
  z-index: 100;
}

.message-input-form {
  padding: 15px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 5px;
  border-radius: 50%;
  transition: background 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-btn:hover {
  background: #f0f0f0;
}

.input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.file-preview {
  font-size: 0.75rem;
  color: var(--color-gold);
  margin-top: 2px;
  padding-left: 10px;
}

.no-chat-selected {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #888;
}

.no-chat-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  opacity: 0.2;
}
</style>
