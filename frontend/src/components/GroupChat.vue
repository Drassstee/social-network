<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useChatStore } from '../stores/chat'
import { useAuthStore } from '../stores/auth'

const props = defineProps(['groupId'])
const chatStore = useChatStore()
const authStore = useAuthStore()

const messages = ref([])
const newMessage = ref('')
const messageContainer = ref(null)

const fetchMessages = async () => {
    try {
        const res = await fetch(`/api/v1/groups/${props.groupId}/messages`)
        const data = await res.json()
        messages.value = data.reverse() || []
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

const handleSendMessage = () => {
    if (!newMessage.value.trim()) return
    
    chatStore.sendGroupMessage(parseInt(props.groupId), newMessage.value)
    
    // Optimistic UI
    messages.value.push({
        id: Date.now(),
        group_id: props.groupId,
        sender_id: authStore.user.id,
        username: authStore.user.nickname || authStore.user.first_name,
        body: newMessage.value,
        created_at: new Date().toISOString()
    })
    
    newMessage.value = ''
    scrollToBottom()
}

// Listen to socket messages
onMounted(() => {
    fetchMessages()
    
    // Subscribe to chat store status or messages
    const unsubscribe = chatStore.$onAction(({ name, args }) => {
        if (name === 'addGroupMessage') {
            const msg = args[0]
            if (msg.group_id === parseInt(props.groupId)) {
                // Check if already added (to avoid double entry from optimistic UI)
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
    <div class="group-chat">
        <div class="messages-container" ref="messageContainer">
            <div v-for="msg in messages" :key="msg.id" 
                 :class="['message', { 'own': msg.sender_id === authStore.user?.id }]">
                <div class="msg-bubble">
                    <div class="msg-header">
                        <strong>{{ msg.username }}</strong>
                        <span>{{ new Date(msg.created_at).toLocaleTimeString() }}</span>
                    </div>
                    <p>{{ msg.body }}</p>
                </div>
            </div>
        </div>

        <div class="input-area">
            <input v-model="newMessage" @keyup.enter="handleSendMessage" placeholder="Send a message..." />
            <button @click="handleSendMessage" class="btn btn-primary">Send</button>
        </div>
    </div>
</template>

<style scoped>
.group-chat {
    height: 500px;
    display: flex;
    flex-direction: column;
    background: var(--color-paper);
    border-radius: var(--border-radius);
    border: 1px solid var(--color-gold);
}

.messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 15px;
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
    background: #f0f0f0;
    padding: 10px 15px;
    border-radius: 12px;
    max-width: 80%;
}

.message.own .msg-bubble {
    background: var(--color-gold);
    color: var(--color-charcoal);
}

.msg-header {
    display: flex;
    justify-content: space-between;
    gap: 10px;
    font-size: 0.75rem;
    margin-bottom: 5px;
    opacity: 0.7;
}

.input-area {
    padding: 15px;
    border-top: 1px solid var(--color-gold);
    display: flex;
    gap: 10px;
}

.input-area input {
    flex: 1;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: var(--border-radius);
    background: transparent;
}
</style>
