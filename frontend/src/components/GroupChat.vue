<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'

const props = defineProps(['groupId'])
const chatStore = useChatStore()
const authStore = useAuthStore()

const messages = ref([])
const newMessage = ref('')
const messageContainer = ref(null)
const showEmojiPicker = ref(false)
const imageFile = ref(null)

const fetchMessages = async () => {
    try {
        const res = await fetch(`/api/v1/groups/${props.groupId}/messages`)
        const data = await res.json()
        messages.value = data || []
        scrollToBottom()
    } catch (err) {
        console.error('Failed to fetch group messages:', err)
    }
}

const scrollToBottom = async () => {
    await nextTick()
    if (messageContainer.value) {
        messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
}

const handleFileChange = (e) => {
    imageFile.value = e.target.files[0]
}

const onSelectEmoji = (emoji) => {
    newMessage.value += emoji.i
}

const handleSendMessage = async () => {
    if (!newMessage.value.trim() && !imageFile.value) return
    
    let imageURL = null
    if (imageFile.value) {
        const formData = new FormData()
        formData.append('image', imageFile.value)
        try {
            const resp = await fetch('/api/v1/chat/upload', { method: 'POST', body: formData })
            if (resp.ok) {
                const data = await resp.json()
                imageURL = data.url
            }
        } catch (err) { console.error('Upload failed:', err) }
    }

    chatStore.sendGroupMessage(parseInt(props.groupId), newMessage.value, imageURL)
    
    // Reset inputs
    newMessage.value = ''
    imageFile.value = null
    showEmojiPicker.value = false
    scrollToBottom()
}

onMounted(() => {
    fetchMessages()
    
    const unsubscribe = chatStore.$onAction(({ name, args }) => {
        if (name === 'addGroupMessage') {
            const msg = args[0]
            if (msg.group_id === parseInt(props.groupId)) {
                if (!messages.value.some(m => m.id === msg.id)) {
                    messages.value.push(msg)
                    scrollToBottom()
                }
            }
        }
    })
    
    onUnmounted(unsubscribe)
})
</script>

<template>
    <div class="group-chat card-traditional">
        <div class="messages-container" ref="messageContainer">
            <div v-for="msg in messages" :key="msg.id" 
                 :class="['message', { 'own': msg.sender_id === authStore.user?.id }]">
                <div class="msg-bubble">
                    <div class="msg-header">
                        <strong>{{ msg.username || msg.sender_nickname }}</strong>
                        <span>{{ new Date(msg.created_at).toLocaleTimeString() }}</span>
                    </div>
                    <img v-if="msg.image_url" :src="msg.image_url" class="chat-img" />
                    <p v-if="msg.body">{{ msg.body }}</p>
                </div>
            </div>
        </div>

        <div class="input-controls">
            <div v-if="showEmojiPicker" class="emoji-picker-container">
                <EmojiPicker :native="true" @select="onSelectEmoji" />
            </div>

            <div class="input-area">
                <button type="button" @click="showEmojiPicker = !showEmojiPicker" class="icon-btn">😊</button>
                <label class="icon-btn">
                    🖼️
                    <input type="file" @change="handleFileChange" accept="image/*" hidden />
                </label>
                <div class="input-wrapper">
                    <input 
                        v-model="newMessage" 
                        @keyup.enter="handleSendMessage" 
                        placeholder="Type a group message..." 
                        class="input-traditional"
                        @focus="showEmojiPicker = false"
                    />
                    <div v-if="imageFile" class="file-preview">📎 {{ imageFile.name }}</div>
                </div>
                <button @click="handleSendMessage" class="btn-traditional mini">Send</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.group-chat {
    height: 600px;
    display: flex;
    flex-direction: column;
    background: var(--color-washi-white);
    padding: 0;
    overflow: hidden;
    border: 1px solid rgba(0,0,0,0.1);
}

.messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 15px;
    scroll-behavior: smooth;
}

.message {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

.message.own {
    align-items: flex-end;
}

.msg-bubble {
    background: white;
    padding: 12px 15px;
    border-radius: 15px;
    max-width: 80%;
    box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.message.own .msg-bubble {
    background: var(--color-gold);
    color: var(--color-charcoal);
}

.msg-header {
    display: flex;
    justify-content: space-between;
    gap: 15px;
    font-size: 0.7rem;
    margin-bottom: 5px;
    opacity: 0.6;
}

.chat-img {
    max-width: 100%;
    border-radius: 8px;
    margin-bottom: 8px;
}

.input-controls {
    position: relative;
    border-top: 1px solid #eee;
    background: white;
}

.emoji-picker-container {
    position: absolute;
    bottom: 100%;
    left: 20px;
    z-index: 100;
}

.input-area {
    padding: 15px;
    display: flex;
    align-items: center;
    gap: 10px;
}

.input-wrapper {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.file-preview {
    font-size: 0.7rem;
    color: var(--color-gold);
    margin-top: 2px;
}

.icon-btn {
    background: none;
    border: none;
    font-size: 1.3rem;
    cursor: pointer;
    padding: 5px;
    border-radius: 50%;
}

.icon-btn:hover {
    background: #f0f0f0;
}

.btn-traditional.mini {
    padding: 8px 15px;
    font-size: 0.85rem;
}
</style>
