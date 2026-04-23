<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const groups = ref([])
const showCreateModal = ref(false)
const newGroup = ref({
  title: '',
  description: ''
})

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
    alert('Join request sent!')
  } catch (err) {
    console.error('Failed to join group:', err)
  }
}

onMounted(fetchGroups)
</script>

<template>
  <div class="groups-view">
    <header class="view-header">
      <div class="header-content">
        <div>
          <h1 class="view-title">Groups</h1>
          <p class="view-subtitle">Groups</p>
        </div>
        <button @click="showCreateModal = true" class="btn-traditional">Create New Group</button>
      </div>
    </header>

    <div class="groups-grid">
      <div v-for="group in groups" :key="group.id" class="card-traditional group-card">
        <div class="group-header">
          <div class="group-icon">Grp</div>
          <div>
            <h3 class="group-title">{{ group.title }}</h3>
            <span class="member-count">{{ group.members?.length || 0 }} members</span>
          </div>
        </div>
        <p class="group-desc">{{ group.description }}</p>
        <div class="group-actions">
          <button @click="joinGroup(group.id)" class="btn-traditional ghost">Request to Join</button>
          <router-link :to="`/groups/${group.id}`" class="btn-traditional mini">View</router-link>
        </div>
      </div>
    </div>

    <!-- Create Group Modal -->
    <div v-if="showCreateModal" class="modal-overlay">
      <div class="card-traditional modal-card">
        <h2>Create New Group</h2>
        <form @submit.prevent="handleCreateGroup" class="modal-form">
          <div class="form-group">
            <label>Title</label>
            <input v-model="newGroup.title" type="text" class="input-traditional" required />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="newGroup.description" class="input-traditional" rows="4" required></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn-traditional ghost">Cancel</button>
            <button type="submit" class="btn-traditional">Create</button>
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

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 30px;
}

.group-card {
  display: flex;
  flex-direction: column;
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
  background: var(--color-charcoal);
  color: var(--color-gold);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-family: 'Noto Serif JP', serif;
}

.group-title {
  font-size: 1.2rem;
}

.member-count {
  font-size: 0.8rem;
  color: #666;
}

.group-desc {
  flex: 1;
  color: var(--color-charcoal);
  margin-bottom: 20px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.group-actions {
  display: flex;
  gap: 10px;
}

.btn-traditional.ghost {
  background: none;
  color: var(--color-vermilion);
  border: 1px solid var(--color-vermilion);
}

.btn-traditional.mini {
  padding: 5px 15px;
  background: var(--color-charcoal);
  color: var(--color-gold);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-card {
  width: 100%;
  max-width: 500px;
}

.modal-form {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 10px;
}
</style>
