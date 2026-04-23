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
        <button @click="showCreateForm = !showCreateForm" class="btn btn-secondary">
            {{ showCreateForm ? 'Cancel' : 'Create New Event' }}
        </button>

        <div v-if="showCreateForm" class="create-event-form card-traditional">
            <input v-model="newEvent.title" placeholder="Event Title" required />
            <textarea v-model="newEvent.description" placeholder="Description"></textarea>
            <input v-model="newEvent.event_time" type="datetime-local" required />
            <button @click="handleCreateEvent" class="btn btn-primary">Save Event</button>
        </div>

        <div class="events-list">
            <div v-for="event in events" :key="event.id" class="event-card">
                <h3>{{ event.title }}</h3>
                <p>{{ event.description }}</p>
                <div class="event-meta">
                    <span>📅 {{ new Date(event.event_time).toLocaleString() }}</span>
                </div>
                <div class="stats">
                    <span>Going: {{ event.going_count }}</span>
                    <span>Not Going: {{ event.not_going_count }}</span>
                </div>
                <div class="actions">
                    <button @click="respondToEvent(event.id, 'going')" class="btn btn-small btn-success">Going</button>
                    <button @click="respondToEvent(event.id, 'not_going')" class="btn btn-small btn-danger">Not Going</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.create-event-form {
    margin-top: 20px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.create-event-form input, .create-event-form textarea {
    padding: 10px;
    border: 1px solid var(--color-gold);
    background: transparent;
}

.events-list {
    margin-top: 20px;
}

.event-card {
    background: var(--color-paper);
    padding: 20px;
    border-radius: var(--border-radius);
    margin-bottom: 20px;
    border: 1px solid var(--color-gold);
}

.event-meta {
    font-size: 0.9rem;
    color: var(--color-charcoal);
    margin: 10px 0;
}

.stats {
    display: flex;
    gap: 20px;
    margin-bottom: 15px;
    font-weight: bold;
}

.actions {
    display: flex;
    gap: 10px;
}

.btn-small {
    padding: 8px 15px;
    font-size: 0.8rem;
}
</style>
