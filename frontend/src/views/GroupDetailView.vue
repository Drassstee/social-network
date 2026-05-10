<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useGroupStore } from '../stores/groups'
import { usePostStore } from '../stores/posts'
import { useUIStore } from '../stores/ui'
import GroupEvents from '../components/GroupEvents.vue'
import GroupChat from '../components/GroupChat.vue'
import CreatePost from '../components/CreatePost.vue'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const auth = useAuthStore()
const groupStore = useGroupStore()
const postStore = usePostStore()
const ui = useUIStore()

const members = ref([])
const activeTab = ref('posts') // 'posts', 'members', 'events', 'chat'
const error = ref(null)

const isMember = computed(() => {
    if (!auth.user || !members.value || !Array.isArray(members.value)) return false
    return members.value.some(m => Number(m.user_id) === Number(auth.user.id))
})

const isCreator = computed(() => {
    if (!auth.user || !groupStore.currentGroup) return false
    return Number(groupStore.currentGroup.creator_id) === Number(auth.user.id)
})
const joinRequests = ref([])
const inviteUserId = ref('')

const fetchGroupData = async () => {
    error.value = null
    const id = route.params.id
    try {
        await groupStore.fetchGroup(id)
        const membersData = await groupStore.fetchMembers(id)
        members.value = Array.isArray(membersData) ? membersData : []
    } catch (err) {
        console.error('Failed to fetch group data:', err)
        error.value = err.message
    }
}

// Watch for membership/creator status changes to fetch private data
watch([isMember, isCreator], ([newMember, newCreator]) => {
    if (newMember || newCreator) {
        postStore.fetchPosts(route.params.id)
        if (newCreator) fetchJoinRequests()
    }
}, { immediate: true })

const fetchJoinRequests = async () => {
    try {
        joinRequests.value = await groupStore.fetchJoinRequests(route.params.id)
    } catch (err) {
        console.error('Failed to fetch join requests:', err)
    }
}

const handleJoinRequest = async () => {
    try {
        await groupStore.requestJoin(route.params.id)
        ui.showToast('Request sent to creator', 'success')
    } catch (err) {
        console.error('Failed to send request:', err)
        ui.showToast('Failed to send request', 'error')
    }
}


const respondToRequest = async (requestId, accept) => {
    try {
        await groupStore.respondToRequest(requestId, accept)
        fetchJoinRequests()
        fetchGroupData()
    } catch (err) {
        console.error('Failed to respond to request:', err)
    }
}


const handleInvite = async (userId) => {
    try {
        await groupStore.inviteUser(groupStore.currentGroup.id, userId)
        ui.showToast('Invitation sent!', 'success')
        inviteUserId.value = ''
    } catch (err) {
        console.error('Failed to invite user:', err)
        ui.showToast('Failed to invite user', 'error')
    }
}

onMounted(fetchGroupData)
watch(() => route.params.id, fetchGroupData)
</script>

<template>
    <div v-if="groupStore.loading" class="loading-container">
        <div class="glow-text pulse">SYSTEM_SYNC_IN_PROGRESS...</div>
    </div>
    
    <div v-else-if="error" class="error-container">
        <div class="error-box card-retro">
            <h2 class="error-title">CRITICAL_ERROR</h2>
            <p class="error-msg">{{ error }}</p>
            <button @click="fetchGroupData" class="btn-retro">RETRY_SYNC</button>
        </div>
    </div>

    <div class="group-detail" v-else-if="groupStore.currentGroup">
        <header class="group-header">
            <div class="group-info">
                <h1 class="glow-text">{{ groupStore.currentGroup.title }}</h1>
                <p class="description">MISSION_LOG: {{ groupStore.currentGroup.description }}</p>
            </div>
            <div class="group-actions">
                <button v-if="!isMember && !isCreator" @click="handleJoinRequest" class="btn-retro">REQUEST_ENTRY</button>
                <div v-else class="member-badge">{{ isCreator ? 'CREATOR_UNIT' : 'AUTHORIZED_UNIT' }}</div>
            </div>
        </header>

        <nav class="group-tabs">
            <button :class="{ active: activeTab === 'posts' }" @click="activeTab = 'posts'">LOGS</button>
            <button :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">ROSTER</button>
            <button v-if="isMember || isCreator" :class="{ active: activeTab === 'events' }" @click="activeTab = 'events'">OPS</button>
            <button v-if="isMember || isCreator" :class="{ active: activeTab === 'chat' }" @click="activeTab = 'chat'">COMM_LINK</button>
            <button v-if="isCreator && joinRequests.length" :class="{ active: activeTab === 'requests' }" @click="activeTab = 'requests'">
                ENTRY_REQS ({{ joinRequests.length }})
            </button>
        </nav>

        <div class="tab-content">
            <!-- Posts Tab -->
            <div v-if="activeTab === 'posts'" class="posts-tab">
                <CreatePost v-if="isMember || isCreator" :groupId="route.params.id" />

                <div v-if="!isMember && !isCreator" class="blocked-content">
                    <p>ENCRYPTED: ONLY AUTHORIZED UNITS CAN VIEW LOGS.</p>
                </div>
                <div v-else class="posts-list">
                    <PostCard v-for="post in postStore.posts" :key="post.id" :post="post" />
                </div>
            </div>

            <!-- Members Tab -->
            <div v-if="activeTab === 'members'" class="members-tab">
                <div v-if="isMember || isCreator" class="invite-section card-retro">
                    <h3>INVITE_NEW_UNIT</h3>
                    <div class="invite-controls">
                      <input v-model="inviteUserId" placeholder="USER_ID" type="number" class="input-retro" />
                      <button @click="handleInvite(parseInt(inviteUserId))" class="btn-retro">INVITE</button>
                    </div>
                </div>
                <div class="members-list">
                  <div v-for="member in members" :key="member.user_id" class="member-item card-retro">
                      <span class="member-name">{{ member.first_name }} {{ member.last_name }}</span>
                      <span class="member-role">[{{ member.role }}]</span>
                  </div>
                </div>
            </div>

            <!-- Events Tab -->
            <div v-if="activeTab === 'events'" class="events-tab">
                <GroupEvents :groupId="groupStore.currentGroup.id" />
            </div>

            <!-- Chat Tab -->
            <div v-if="activeTab === 'chat'" class="chat-tab">
                <GroupChat :groupId="groupStore.currentGroup.id" />
            </div>

            <!-- Join Requests Tab -->
            <div v-if="activeTab === 'requests'" class="requests-tab">
                <div v-for="req in joinRequests" :key="req.id" class="request-item card-retro">
                    <span>UNIT_{{ req.username }} ACCESS_REQEUST</span>
                    <div class="actions">
                        <button @click="respondToRequest(req.id, true)" class="btn-retro mini success">ALLOW</button>
                        <button @click="respondToRequest(req.id, false)" class="btn-retro mini danger">DENY</button>
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
    background: rgba(31, 11, 53, 0.85);
    padding: 30px;
    border: 2px solid var(--color-neon-cyan);
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    box-shadow: 8px 8px 0 var(--color-neon-magenta);
}

.glow-text {
    color: var(--color-neon-cyan);
    text-shadow: var(--shadow-neon-cyan);
    margin: 0 0 10px 0;
}

.description {
    color: var(--color-neon-yellow);
    font-family: 'VT323', monospace;
    font-size: 1.4rem;
}

.member-badge {
    color: var(--color-neon-yellow);
    border: 1px solid var(--color-neon-yellow);
    padding: 5px 15px;
    font-family: 'VT323', monospace;
    text-shadow: 0 0 5px var(--color-neon-yellow);
}

.group-tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    border-bottom: 2px solid var(--color-neon-magenta);
    padding-bottom: 5px;
    flex-wrap: wrap;
}

.group-tabs button {
    background: none;
    border: none;
    color: var(--color-neon-cyan);
    padding: 10px 20px;
    cursor: pointer;
    font-weight: bold;
    font-family: 'Press Start 2P', cursive;
    font-size: 0.7rem;
    transition: all 0.2s;
}

.group-tabs button:hover {
  color: var(--color-neon-magenta);
  text-shadow: 0 0 5px var(--color-neon-magenta);
}

.group-tabs button.active {
    background: var(--color-neon-cyan);
    color: var(--color-dark-bg);
}

.create-post {
    margin-bottom: 20px;
}

.post-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 15px;
}

.file-label {
    cursor: pointer;
    background: rgba(0, 255, 255, 0.1);
    border: 1px solid var(--color-neon-cyan);
    padding: 8px 15px;
    color: var(--color-neon-cyan);
    font-family: 'VT323', monospace;
    font-size: 1.2rem;
}

.hidden-input { display: none; }

.post-card {
    margin-bottom: 20px;
}

.post-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    font-family: 'VT323', monospace;
}

.author-name {
    color: var(--color-neon-magenta);
    font-size: 1.2rem;
}

.post-date {
    color: var(--color-neon-yellow);
}

.post-content {
    font-family: 'VT323', monospace;
    font-size: 1.3rem;
    color: white;
    background: rgba(0,0,0,0.2);
    padding: 10px;
}

.post-image {
    max-width: 100%;
    max-height: 400px;
    object-fit: contain;
    margin-top: 10px;
    border: 1px solid var(--color-neon-cyan);
    display: block;
}

.invite-section {
  margin-bottom: 20px;
}

.invite-section h3 {
  font-family: 'Press Start 2P', cursive;
  font-size: 0.8rem;
  margin-bottom: 15px;
  color: var(--color-neon-magenta);
}

.invite-controls {
  display: flex;
  gap: 15px;
}

.members-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
}

.member-item {
    display: flex;
    flex-direction: column;
    padding: 15px;
    border-radius: 0;
}

.member-name {
    color: var(--color-neon-cyan);
    font-family: 'VT323', monospace;
    font-size: 1.2rem;
}

.member-role {
    font-family: 'VT323', monospace;
    color: var(--color-neon-yellow);
}

.request-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    margin-bottom: 10px;
}

.request-item span {
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  color: var(--color-neon-yellow);
}

.btn-retro.mini {
    padding: 5px 10px;
    font-size: 0.6rem;
}

.btn-retro.mini.success { border-color: #00ff00; color: #00ff00; }
.btn-retro.mini.danger { border-color: #ff0000; color: #ff0000; }

.blocked-content {
  text-align: center;
  padding: 50px;
  border: 2px dashed var(--color-neon-magenta);
  color: var(--color-neon-magenta);
  font-family: 'Press Start 2P', cursive;
  font-size: 0.8rem;
}
.loading-container, .error-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
    text-align: center;
}

.pulse {
    animation: pulse 1.5s infinite;
    font-family: 'Press Start 2P', cursive;
    font-size: 1rem;
}

@keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.3; }
    100% { opacity: 1; }
}

.error-box {
    padding: 40px;
    border-color: #ff0000;
    box-shadow: 0 0 20px rgba(255, 0, 0, 0.2);
}

.error-title {
    color: #ff0000;
    text-shadow: 0 0 10px #ff0000;
    font-family: 'Press Start 2P', cursive;
    margin-bottom: 20px;
}

.error-msg {
    color: var(--color-neon-yellow);
    font-family: 'VT323', monospace;
    font-size: 1.5rem;
    margin-bottom: 30px;
}
</style>
