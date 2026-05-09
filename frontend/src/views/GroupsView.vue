<script setup>
import { ref, onMounted } from 'vue'
import { useChatStore } from '../stores/chat'

const chat = useChatStore()
const groups = ref([])
const showCreateModal = ref(false)
const newGroup = ref({
  title: '',
  description: ''
})
const feedbackMessage = ref('')
const feedbackType = ref('info')

const showFeedback = (msg, type = 'info') => {
  feedbackMessage.value = msg
  feedbackType.value = type
  setTimeout(() => { feedbackMessage.value = '' }, 3000)
}

const fetchGroups = async () => {
  try {
    const response = await fetch('/api/v1/groups')
    const data = await response.json()
    groups.value = data || []
  } catch (err) {
    console.error('Failed to fetch groups:', err)
  }
}

const handleCreateGroup = async () => {
  try {
    const response = await fetch('/api/v1/groups', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newGroup.value)
    })
    if (response.ok) {
      showCreateModal.value = false
      newGroup.value = { title: '', description: '' }
      fetchGroups()
    }
  } catch (err) {
    console.error('Failed to create group:', err)
  }
}

const joinGroup = async (groupId) => {
  try {
    await fetch(`/api/v1/groups/${groupId}/request`, { method: 'POST' })
    showFeedback('Join request sent!', 'success')
  } catch (err) {
    console.error('Failed to join group:', err)
  }
}

onMounted(fetchGroups)
</script>

<template>
  <div class="groups-view">
    <div v-if="feedbackMessage" class="feedback-toast" :class="feedbackType">
      {{ feedbackMessage }}
    </div>
    <header class="view-header">
      <div class="header-content">
        <div>
          <h1 class="view-title">SYNDICATES</h1>
          <p class="view-subtitle">ACTIVE_NODES</p>
        </div>
        <button @click="showCreateModal = true" class="btn-retro">NEW_SYNDICATE</button>
      </div>
    </header>

    <div class="groups-grid">
        <router-link v-for="group in groups" :key="group.id" :to="`/groups/${group.id}`" class="card-retro group-card">
          <div class="group-header">
            <div class="group-icon">#</div>
            <div>
              <h3 class="group-title">{{ group.title }}</h3>
              <span class="member-count">{{ group.member_count || 0 }} UNITS</span>
            </div>
            <div v-if="chat.getUnreadCount(group.id, 'g') > 0" class="badge">{{ chat.getUnreadCount(group.id, 'g') }}</div>
          </div>
          <p class="group-desc">{{ group.description }}</p>
          <div class="group-actions" @click.stop>
            <button @click.prevent="joinGroup(group.id)" class="btn-retro ghost">REQUEST_ACCESS</button>
            <button class="btn-retro mini">OPEN</button>
          </div>
        </router-link>
    </div>

    <!-- Create Group Modal -->
    <div v-if="showCreateModal" class="modal-overlay">
      <div class="card-retro modal-card">
        <h2>NEW_SYNDICATE_CONFIG</h2>
        <form @submit.prevent="handleCreateGroup" class="modal-form">
          <div class="form-group">
            <label>SYNDICATE_TITLE</label>
            <input v-model="newGroup.title" type="text" class="input-retro" required />
          </div>
          <div class="form-group">
            <label>MISSION_PARAMETERS</label>
            <textarea v-model="newGroup.description" class="input-retro" rows="4" required></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-retro ghost">ABORT</button>
            <button type="submit" class="btn-retro">INITIALIZE</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-header {
  margin-bottom: 40px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.view-title {
  font-size: 3rem;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
}

.view-subtitle {
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 30px;
}

.group-card {
  display: flex;
  flex-direction: column;
  transition: all 0.2s;
  text-decoration: none;
}

.group-card:hover {
  transform: translateY(-5px);
  box-shadow: 12px 12px 0 var(--color-neon-cyan);
}

.group-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.group-icon {
  width: 50px;
  height: 50px;
  background: var(--color-dark-bg);
  color: var(--color-neon-magenta);
  border: 2px solid var(--color-neon-magenta);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-family: 'Press Start 2P', cursive;
  box-shadow: 2px 2px 0 var(--color-neon-cyan);
}

.group-title {
  font-size: 1.2rem;
  color: var(--color-neon-cyan);
}

.member-count {
  font-size: 0.9rem;
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
}

.group-desc {
  flex: 1;
  color: white;
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  margin-bottom: 20px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  background: rgba(0,0,0,0.2);
  padding: 10px;
}

.group-actions {
  display: flex;
  gap: 10px;
}

.btn-retro.ghost {
  background: rgba(255, 0, 255, 0.1);
  font-size: 0.7rem;
}

.btn-retro.mini {
  padding: 5px 15px;
  font-size: 0.7rem;
  box-shadow: 2px 2px 0 var(--color-neon-magenta);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(11, 12, 16, 0.85);
  backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-card {
  width: 100%;
  max-width: 500px;
}

.modal-card h2 {
  font-size: 1.5rem;
  color: var(--color-neon-magenta);
  text-shadow: var(--shadow-neon-magenta);
}

.modal-form {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group label {
  color: var(--color-neon-cyan);
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  margin-bottom: 5px;
  display: block;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 10px;
}
</style>
