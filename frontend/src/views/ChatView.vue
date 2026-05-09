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
    try {
      imageURL = await chat.uploadImage(imageFile.value)
    } catch (err) { 
      console.error('Upload failed:', err)
      return // Don't send message if upload fails
    }
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
  <div class="chat-container card-retro">
    <div class="users-sidebar">
      <h3 class="sidebar-title">COMM_UNITS</h3>
      <div class="users-list">
        <div 
          v-for="user in chat.onlineUsers" 
          :key="user.id"
          class="user-item"
          :class="{ active: chat.activeChatUser?.id === user.id }"
          @click="selectUser(user)"
        >
          <div class="avatar-sm">{{ user.username?.[0] || '?' }}</div>
          <div class="user-info">
            <span class="username">{{ user.username }}</span>
            <span class="status-text" :class="{ online: user.is_online }">
              {{ user.is_online ? 'ONLINE' : 'OFFLINE' }}
            </span>
          </div>
          <div v-if="chat.getUnreadCount(user.id) > 0" class="badge mini">{{ chat.getUnreadCount(user.id) }}</div>
          <div class="online-indicator" :class="{ offline: !user.is_online }"></div>
        </div>
      </div>
    </div>

    
    <div class="chat-main">
      <template v-if="chat.activeChatUser">
        <div class="chat-header">
          <h3 class="glow-text">COMM_LINK: {{ chat.activeChatUser.username }}</h3>
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
            <button type="button" @click="showEmojiPicker = !showEmojiPicker" class="icon-btn">💾</button>
            <label class="icon-btn file-label">
              📎
              <input type="file" @change="handleFileChange" accept="image/*" hidden />
            </label>
            <div class="input-wrapper">
              <input 
                v-model="newMessage" 
                type="text" 
                placeholder="INPUT_MESSAGE..." 
                class="input-retro"
                @focus="showEmojiPicker = false"
              />
              <div v-if="imageFile" class="file-preview">ATTACHED: {{ imageFile.name }}</div>
            </div>
            <button type="submit" class="btn-retro">TRANSMIT</button>
          </form>
        </div>
      </template>
      <div v-else class="no-chat-selected">
        <div class="no-chat-icon">>_</div>
        <p>SELECT_UNIT_TO_ESTABLISH_COMM_LINK</p>
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
  background: var(--color-dark-bg);
  border: 2px solid var(--color-neon-cyan);
  box-shadow: 10px 10px 0 var(--color-neon-magenta);
}

.users-sidebar {
  width: 250px;
  border-right: 2px solid var(--color-grid-line);
  display: flex;
  flex-direction: column;
  background: rgba(0,0,0,0.3);
}

.sidebar-title {
  padding: 20px;
  border-bottom: 2px solid var(--color-grid-line);
  font-size: 0.8rem;
  font-family: 'Press Start 2P', cursive;
  background: rgba(31, 11, 53, 0.8);
  color: var(--color-neon-cyan);
  text-shadow: 0 0 5px var(--color-neon-cyan);
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
  border-bottom: 1px solid var(--color-grid-line);
  transition: all 0.2s;
  color: white;
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
}

.user-item:hover {
  background: rgba(255, 0, 255, 0.1);
  color: var(--color-neon-magenta);
}

.user-item.active {
  background: rgba(0, 255, 255, 0.1);
  color: var(--color-neon-cyan);
  border-left: 4px solid var(--color-neon-cyan);
}

.avatar-sm {
  width: 35px;
  height: 35px;
  background: var(--color-dark-bg);
  color: var(--color-neon-magenta);
  border: 1px solid var(--color-neon-magenta);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  box-shadow: 2px 2px 0 var(--color-neon-cyan);
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.status-text {
  font-size: 0.7rem;
  color: #666;
}

.status-text.online {
  color: #00ff00;
  text-shadow: 0 0 2px #00ff00;
}

.online-indicator {
  width: 10px;
  height: 10px;
  background: #00ff00;
  border-radius: 50%;
  margin-left: auto;
  box-shadow: 0 0 8px #00ff00;
}

.online-indicator.offline {
  background: #444;
  box-shadow: none;
}


.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: rgba(11, 12, 16, 0.5);
}

.chat-header {
  padding: 15px 25px;
  border-bottom: 2px solid var(--color-grid-line);
  background: rgba(31, 11, 53, 0.5);
}

.glow-text {
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}

.messages-area {
  flex: 1;
  padding: 25px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
  scroll-behavior: smooth;
  background-image: 
    linear-gradient(var(--color-grid-line) 1px, transparent 1px),
    linear-gradient(90deg, var(--color-grid-line) 1px, transparent 1px);
  background-size: 40px 40px;
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
  background: rgba(31, 11, 53, 0.9);
  border: 1px solid var(--color-neon-magenta);
  position: relative;
  font-family: 'VT323', monospace;
  font-size: 1.3rem;
  color: white;
  box-shadow: 4px 4px 0 rgba(255, 0, 255, 0.3);
}

.sent .message-bubble {
  border-color: var(--color-neon-cyan);
  box-shadow: -4px 4px 0 rgba(0, 255, 255, 0.3);
}

.chat-img {
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
  border: 1px solid var(--color-neon-cyan);
  margin-bottom: 8px;
  display: block;
}

.time {
  font-size: 0.8rem;
  color: var(--color-neon-yellow);
  display: block;
  margin-top: 5px;
  text-align: right;
}

.input-controls {
  position: relative;
  border-top: 2px solid var(--color-grid-line);
  background: rgba(31, 11, 53, 0.8);
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
  gap: 15px;
}

.icon-btn {
  background: none;
  border: 1px solid var(--color-neon-cyan);
  color: var(--color-neon-cyan);
  font-size: 1.2rem;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: rgba(0, 255, 255, 0.1);
  box-shadow: 0 0 5px var(--color-neon-cyan);
}

.input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.file-preview {
  font-size: 0.8rem;
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
  margin-top: 5px;
}

.no-chat-selected {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-neon-magenta);
  font-family: 'Press Start 2P', cursive;
  text-align: center;
  padding: 20px;
}

.no-chat-icon {
  font-size: 5rem;
  margin-bottom: 20px;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

.no-chat-selected p {
  font-size: 0.8rem;
  line-height: 1.5;
}
</style>
