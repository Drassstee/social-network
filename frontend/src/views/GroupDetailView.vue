<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import GroupEvents from '../components/GroupEvents.vue'
import GroupChat from '../components/GroupChat.vue'

const route = useRoute()
const auth = useAuthStore()

const group = ref(null)
const members = ref([])
const posts = ref([])
const activeTab = ref('posts') // 'posts', 'members', 'events', 'chat'
const isLoading = ref(true)
const isMember = ref(false)
const isCreator = ref(false)
const joinRequests = ref([])
const inviteUserId = ref('')
const feedbackMessage = ref('')
const feedbackType = ref('info') // 'success', 'error', 'info'

const showFeedback = (msg, type = 'info') => {
    feedbackMessage.value = msg
    feedbackType.value = type
    setTimeout(() => {
        feedbackMessage.value = ''
    }, 3000)
}

const newPostContent = ref('')
const newPostImage = ref(null)
const isPosting = ref(false)

const fetchGroupData = async () => {
    isLoading.value = true
    try {
        const [groupRes, membersRes] = await Promise.all([
            fetch(`/api/v1/groups/${route.params.id}`),
            fetch(`/api/v1/groups/${route.params.id}/members`)
        ])
        
        group.value = await groupRes.json()
        members.value = await membersRes.json()
        
        isMember.value = members.value.some(m => Number(m.user_id) === Number(auth.user?.id))
        isCreator.value = Number(group.value.creator_id) === Number(auth.user?.id)
        
        if (isMember.value) {
            fetchPosts()
            if (isCreator.value) fetchJoinRequests()
        }
    } catch (err) {
        console.error('Failed to fetch group data:', err)
    } finally {
        isLoading.value = false
    }
}

const fetchPosts = async () => {
    try {
        const res = await fetch(`/api/v1/posts?group_id=${route.params.id}`)
        const data = await res.json()
        posts.value = data.posts || []
    } catch (err) {
        console.error('Failed to fetch group posts:', err)
    }
}

const fetchJoinRequests = async () => {
    try {
        const res = await fetch(`/api/v1/groups/${route.params.id}/requests`)
        joinRequests.value = await res.json()
    } catch (err) {
        console.error('Failed to fetch join requests:', err)
    }
}

const handleJoinRequest = async () => {
    try {
        await fetch(`/api/v1/groups/${route.params.id}/request`, { method: 'POST' })
        showFeedback('Request sent to creator', 'success')
    } catch (err) {
        console.error('Failed to send request:', err)
        showFeedback('Failed to send request', 'error')
    }
}

const handleCreatePost = async () => {
    if (!newPostContent.value && !newPostImage.value) return
    
    isPosting.value = true
    const formData = new FormData()
    formData.append('content', newPostContent.value)
    formData.append('group_id', route.params.id)
    if (newPostImage.value) formData.append('image', newPostImage.value)
    
    try {
        const res = await fetch('/api/v1/posts', {
            method: 'POST',
            body: formData
        })
        if (res.ok) {
            newPostContent.value = ''
            newPostImage.value = null
            fetchPosts()
        }
    } catch (err) {
        console.error('Failed to create post:', err)
    } finally {
        isPosting.value = false
    }
}

const respondToRequest = async (requestId, accept) => {
    try {
        const res = await fetch(`/api/v1/groups/requests/${requestId}/respond`, {
            method: 'POST',
            body: JSON.stringify({ accept })
        })
        if (res.ok) {
            fetchJoinRequests()
            fetchGroupData()
        }
    } catch (err) {
        console.error('Failed to respond to request:', err)
    }
}

const onFileChange = (e) => {
    newPostImage.value = e.target.files[0]
}



const handleInvite = async (userId) => {
    try {
        const res = await fetch(`/api/v1/groups/${group.value.id}/invite`, {
            method: 'POST',
            body: JSON.stringify({ user_id: userId })
        })
        if (res.ok) {
            showFeedback('Invitation sent!', 'success')
            inviteUserId.value = ''
        }
    } catch (err) {
        console.error('Failed to invite user:', err)
        showFeedback('Failed to invite user', 'error')
    }
}

onMounted(fetchGroupData)
watch(() => route.params.id, fetchGroupData)
</script>

<template>
    <div class="group-detail" v-if="group">
        <div v-if="feedbackMessage" class="feedback-toast" :class="feedbackType">
            {{ feedbackMessage }}
        </div>
        <header class="group-header">
            <div class="group-info">
                <h1>{{ group.title }}</h1>
                <p class="description">{{ group.description }}</p>
            </div>
            <div class="group-actions">
                <button v-if="!isMember" @click="handleJoinRequest" class="btn btn-primary">Join Group</button>
                <div v-else class="member-badge">Member</div>
            </div>
        </header>

        <nav class="group-tabs">
            <button :class="{ active: activeTab === 'posts' }" @click="activeTab = 'posts'">Posts</button>
            <button :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">Members</button>
            <button v-if="isMember" :class="{ active: activeTab === 'events' }" @click="activeTab = 'events'">Events</button>
            <button v-if="isMember" :class="{ active: activeTab === 'chat' }" @click="activeTab = 'chat'">Chat</button>
            <button v-if="isCreator && joinRequests.length" :class="{ active: activeTab === 'requests' }" @click="activeTab = 'requests'">
                Requests ({{ joinRequests.length }})
            </button>
        </nav>

        <div class="tab-content">
            <!-- Posts Tab -->
            <div v-if="activeTab === 'posts'" class="posts-tab">
                <div v-if="isMember" class="create-post">
                    <textarea v-model="newPostContent" placeholder="Share something with the group..."></textarea>
                    <div class="post-controls">
                        <input type="file" @change="onFileChange" accept="image/*" />
                        <button @click="handleCreatePost" :disabled="isPosting" class="btn btn-primary">Post</button>
                    </div>
                </div>

                <div v-if="!isMember" class="blocked-content">
                    <p>Only members can see group posts.</p>
                </div>
                <div v-else class="posts-list">
                    <div v-for="post in posts" :key="post.id" class="post-card">
                        <div class="post-header">
                            <strong>{{ post.author.first_name }} {{ post.author.last_name }}</strong>
                            <span>{{ new Date(post.created_at).toLocaleString() }}</span>
                        </div>
                        <p class="post-content">{{ post.content }}</p>
                        <img v-if="post.image_url" :src="post.image_url" class="post-image" />
                    </div>
                </div>
            </div>

            <!-- Members Tab -->
            <div v-if="activeTab === 'members'" class="members-tab">
                <div v-if="isMember" class="invite-section">
                    <input v-model="inviteUserId" placeholder="User ID to invite" type="number" />
                    <button @click="handleInvite(parseInt(inviteUserId))" class="btn btn-secondary">Invite</button>
                    <p class="hint">Tip: Enter User ID to invite them to this group.</p>
                </div>
                <div v-for="member in members" :key="member.user_id" class="member-item">
                    <span class="member-name">{{ member.first_name }} {{ member.last_name }}</span>
                    <span class="member-role">{{ member.role }}</span>
                </div>
            </div>

            <!-- Events Tab -->
            <div v-if="activeTab === 'events'" class="events-tab">
                <GroupEvents :groupId="group.id" />
            </div>

            <!-- Chat Tab -->
            <div v-if="activeTab === 'chat'" class="chat-tab">
                <GroupChat :groupId="group.id" />
            </div>

            <!-- Join Requests Tab -->
            <div v-if="activeTab === 'requests'" class="requests-tab">
                <div v-for="req in joinRequests" :key="req.id" class="request-item">
                    <span>{{ req.username }} wants to join</span>
                    <div class="actions">
                        <button @click="respondToRequest(req.id, true)" class="btn btn-small btn-success">Accept</button>
                        <button @click="respondToRequest(req.id, false)" class="btn btn-small btn-danger">Decline</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.group-detail {
    max-width: 800px;
    margin: 0 auto;
}

.group-header {
    background: var(--color-charcoal);
    padding: 30px;
    border-radius: var(--border-radius);
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    border-bottom: 4px solid var(--color-gold);
}

.group-info h1 {
    color: var(--color-gold);
    margin: 0 0 10px 0;
}

.description {
    color: var(--color-paper);
    opacity: 0.8;
}

.group-tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    border-bottom: 2px solid var(--color-gold);
    padding-bottom: 10px;
}

.group-tabs button {
    background: none;
    border: none;
    color: var(--color-charcoal);
    padding: 10px 20px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.3s;
}

.group-tabs button.active {
    background: var(--color-gold);
    color: var(--color-charcoal);
    border-radius: var(--border-radius);
}

.create-post {
    background: var(--color-paper);
    padding: 20px;
    border-radius: var(--border-radius);
    margin-bottom: 20px;
}

.create-post textarea {
    width: 100%;
    min-height: 100px;
    background: transparent;
    border: 1px solid var(--color-gold);
    padding: 10px;
    color: var(--color-charcoal);
    margin-bottom: 10px;
}

.post-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.post-card {
    background: var(--color-paper);
    padding: 20px;
    border-radius: var(--border-radius);
    margin-bottom: 15px;
    border: 1px solid rgba(0,0,0,0.1);
}

.post-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    font-size: 0.9rem;
}

.post-image {
    max-width: 100%;
    margin-top: 10px;
    border-radius: var(--border-radius);
}

.member-item, .request-item {
    display: flex;
    justify-content: space-between;
    background: var(--color-paper);
    padding: 15px;
    border-radius: var(--border-radius);
    margin-bottom: 10px;
}

.member-role {
    font-style: italic;
    color: var(--color-gold);
}

.btn-small {
    padding: 5px 10px;
    margin-left: 5px;
}

.btn-success { background: #4CAF50; color: white; }
.btn-danger { background: #f44336; color: white; }
</style>
