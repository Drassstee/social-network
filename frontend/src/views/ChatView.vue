<script setup>
import { ref, computed } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'

const chat = useChatStore()
const auth = useAuthStore()
const newMessage = ref('')

const activeChatMessages = computed(() => {
  if (!chat.activeChatUser) return []
  return chat.messages[chat.activeChatUser.id] || []
})

const selectUser = (user) => {
  chat.activeChatUser = user
}

const handleSendMessage = () => {
  if (!newMessage.value.trim() || !chat.activeChatUser) return
  chat.sendMessage(chat.activeChatUser.id, newMessage.value)
  newMessage.value = ''
}
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
        
        <div class="messages-area">
          <div 
            v-for="msg in activeChatMessages" 
            :key="msg.id"
            class="message-wrapper"
            :class="{ sent: msg.sender_id === auth.user.id }"
          >
            <div class="message-bubble">
              <p>{{ msg.body }}</p>
              <span class="time">{{ new Date(msg.created_at).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}) }}</span>
            </div>
          </div>
        </div>
        
        <form @submit.prevent="handleSendMessage" class="message-input-form">
          <input 
            v-model="newMessage" 
            type="text" 
            placeholder="Type a message..." 
            class="input-traditional"
          />
          <button type="submit" class="btn-traditional">Send</button>
        </form>
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
  height: 70vh;
  padding: 0;
  overflow: hidden;
}

.users-sidebar {
  width: 250px;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
}

.sidebar-title {
  padding: 20px;
  border-bottom: 1px solid #eee;
  font-size: 1.1rem;
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
  background: var(--color-paper);
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
  font-size: 0.8rem;
}

.online-indicator {
  width: 10px;
  height: 10px;
  background: #4caf50;
  border-radius: 50%;
  margin-left: auto;
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
}

.messages-area {
  flex: 1;
  padding: 25px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 15px;
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
  background: #eee;
  position: relative;
}

.sent .message-bubble {
  background: var(--color-vermilion);
  color: white;
}

.time {
  font-size: 0.7rem;
  opacity: 0.7;
  display: block;
  margin-top: 5px;
  text-align: right;
}

.message-input-form {
  padding: 20px;
  display: flex;
  gap: 15px;
  border-top: 1px solid #eee;
  background: white;
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
  opacity: 0.3;
}
</style>
