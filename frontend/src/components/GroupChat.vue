<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'

const props = defineProps(['groupId'])
const chatStore = useChatStore()
const authStore = useAuthStore()

const newMessage = ref('')
const messageContainer = ref(null)
const showEmojiPicker = ref(false)
const imageFile = ref(null)

const fetchMessages = async () => {
    chatStore.activeGroupID = parseInt(props.groupId)
    await chatStore.fetchGroupMessages(props.groupId)
    scrollToBottom()
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
        try {
            imageURL = await chatStore.uploadImage(imageFile.value)
        } catch (err) { 
            console.error('Upload failed:', err)
            return
        }
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
    
    // Auto-scroll on new messages
    watch(() => chatStore.groupMessages[props.groupId], () => {
        scrollToBottom()
    }, { deep: true })
    
    onUnmounted(() => {
        chatStore.activeGroupID = null
    })
})
</script>

<template>
    <div class="group-chat card-retro">
        <div class="messages-container" ref="messageContainer">
            <div v-for="msg in (chatStore.groupMessages[props.groupId] || [])" :key="msg.id" 
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
                <button type="button" @click="showEmojiPicker = !showEmojiPicker" class="icon-btn">💾</button>
                <label class="icon-btn">
                    📎
                    <input type="file" @change="handleFileChange" accept="image/*" hidden />
                </label>
                <div class="input-wrapper">
                    <input 
                        v-model="newMessage" 
                        @keyup.enter="handleSendMessage" 
                        placeholder="SYNDICATE_COMMS..." 
                        class="input-retro"
                        @focus="showEmojiPicker = false"
                    />
                    <div v-if="imageFile" class="file-preview">ATTACHED: {{ imageFile.name }}</div>
                </div>
                <button @click="handleSendMessage" class="btn-retro mini">TRANSMIT</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.group-chat {
    height: 600px;
    display: flex;
    flex-direction: column;
    background: rgba(11, 12, 16, 0.5);
    padding: 0;
    overflow: hidden;
    border: 2px solid var(--color-neon-magenta);
}

.messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    scroll-behavior: smooth;
    background-image: 
      linear-gradient(rgba(255, 0, 255, 0.05) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 0, 255, 0.05) 1px, transparent 1px);
    background-size: 30px 30px;
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
    background: rgba(31, 11, 53, 0.9);
    padding: 12px 18px;
    border: 1px solid var(--color-neon-magenta);
    max-width: 80%;
    box-shadow: 4px 4px 0 rgba(255, 0, 255, 0.2);
    font-family: 'VT323', monospace;
    font-size: 1.3rem;
    color: white;
}

.message.own .msg-bubble {
    border-color: var(--color-neon-cyan);
    box-shadow: -4px 4px 0 rgba(0, 255, 255, 0.2);
}

.msg-header {
    display: flex;
    justify-content: space-between;
    gap: 15px;
    font-size: 0.8rem;
    margin-bottom: 5px;
    color: var(--color-neon-yellow);
}

.chat-img {
    max-width: 100%;
    max-height: 250px;
    object-fit: contain;
    border: 1px solid var(--color-neon-cyan);
    margin-bottom: 8px;
    display: block;
}

.input-controls {
    position: relative;
    border-top: 2px solid var(--color-grid-line);
    background: rgba(31, 11, 53, 0.8);
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
    gap: 12px;
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

.btn-retro.mini {
    padding: 8px 15px;
    font-size: 0.7rem;
}
</style>
