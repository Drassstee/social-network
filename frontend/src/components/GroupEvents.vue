<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps(['groupId'])
const events = ref([])
const showCreateForm = ref(false)
const newEvent = ref({
    title: '',
    description: '',
    event_time: ''
})

const fetchEvents = async () => {
    try {
        const res = await fetch(`/api/v1/groups/${props.groupId}/events`)
        events.value = await res.json()
    } catch (err) {
        console.error('Failed to fetch events:', err)
    }
}

const handleCreateEvent = async () => {
    try {
        const res = await fetch(`/api/v1/groups/${props.groupId}/events`, {
            method: 'POST',
            body: JSON.stringify({
                ...newEvent.value,
                event_time: new Date(newEvent.value.event_time).toISOString()
            })
        })
        if (res.ok) {
            showCreateForm.value = false
            newEvent.value = { title: '', description: '', event_time: '' }
            fetchEvents()
        }
    } catch (err) {
        console.error('Failed to create event:', err)
    }
}

const respondToEvent = async (eventId, response) => {
    try {
        await fetch(`/api/v1/groups/events/${eventId}/respond`, {
            method: 'POST',
            body: JSON.stringify({ response })
        })
        fetchEvents()
    } catch (err) {
        console.error('Failed to respond to event:', err)
    }
}

onMounted(fetchEvents)
</script>

<template>
    <div class="group-events">
        <button @click="showCreateForm = !showCreateForm" class="btn-retro">
            {{ showCreateForm ? 'ABORT_CONFIG' : 'NEW_MISSION_OBJECTIVE' }}
        </button>

        <div v-if="showCreateForm" class="create-event-form card-retro">
            <h3 class="form-title">EVENT_PARAMETERS</h3>
            <input v-model="newEvent.title" class="input-retro" placeholder="MISSION_TITLE" required />
            <textarea v-model="newEvent.description" class="input-retro" placeholder="BRIEFING_DETAILS"></textarea>
            <input v-model="newEvent.event_time" type="datetime-local" class="input-retro" required />
            <button @click="handleCreateEvent" class="btn-retro">SAVE_OBJECTIVE</button>
        </div>

        <div class="events-list">
            <div v-for="event in events" :key="event.id" class="event-card card-retro">
                <h3 class="glow-text">{{ event.title }}</h3>
                <p class="event-desc">{{ event.description }}</p>
                <div class="event-meta">
                    <span>TIMELOCK: {{ new Date(event.event_time).toLocaleString() }}</span>
                </div>
                <div class="stats">
                    <span class="stat-item">UNITS_CONFIRMED: {{ event.going_count }}</span>
                    <span class="stat-item">UNITS_REJECTED: {{ event.not_going_count }}</span>
                </div>
                <div class="actions">
                    <button @click="respondToEvent(event.id, 'going')" class="btn-retro mini success">CONFIRM</button>
                    <button @click="respondToEvent(event.id, 'not_going')" class="btn-retro mini danger">REJECT</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.group-events {
  margin-top: 20px;
}

.create-event-form {
    margin-top: 20px;
    padding: 25px;
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.form-title {
  font-family: 'Press Start 2P', cursive;
  font-size: 0.8rem;
  color: var(--color-neon-magenta);
  margin-bottom: 10px;
}

.events-list {
    margin-top: 30px;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
}

.event-card {
    padding: 20px;
    display: flex;
    flex-direction: column;
}

.glow-text {
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
  font-size: 1.5rem;
  margin-bottom: 10px;
}

.event-desc {
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  color: white;
  margin-bottom: 15px;
  flex: 1;
}

.event-meta {
    font-size: 1rem;
    color: var(--color-neon-yellow);
    font-family: 'VT323', monospace;
    margin-bottom: 15px;
}

.stats {
    display: flex;
    flex-direction: column;
    gap: 5px;
    margin-bottom: 20px;
    font-family: 'VT323', monospace;
}

.stat-item {
  font-size: 1.1rem;
  color: #aaa;
}

.actions {
    display: flex;
    gap: 15px;
}

.btn-retro.mini {
    padding: 8px 15px;
    font-size: 0.7rem;
}

.btn-retro.mini.success { border-color: #00ff00; color: #00ff00; }
.btn-retro.mini.danger { border-color: #ff0000; color: #ff0000; }
</style>
