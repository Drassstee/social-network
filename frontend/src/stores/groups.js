import { defineStore } from 'pinia'
import { api } from '../api/api'

export const useGroupStore = defineStore('groups', {
  state: () => ({
    groups: [],
    currentGroup: null,
    loading: false,
    error: null
  }),
  actions: {
    async fetchGroups() {
      this.loading = true
      try {
        const data = await api.get('/groups')
        this.groups = data || []
      } catch (err) {
        this.error = err.message
        console.error('Failed to fetch groups:', err)
      } finally {
        this.loading = false
      }
    },
    async fetchGroup(groupId) {
        this.loading = true
        try {
            const data = await api.get(`/groups/${groupId}`)
            this.currentGroup = data
            return data
        } catch (err) {
            this.error = err.message
            throw err
        } finally {
            this.loading = false
        }
    },
    async createGroup(title, description) {
      try {
        const data = await api.post('/groups', { title, description })
        this.groups.push(data)
        return data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async requestJoin(groupId) {
      try {
        await api.post(`/groups/${groupId}/request`)
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async inviteUser(groupId, userId) {
        try {
            await api.post(`/groups/${groupId}/invite`, { user_id: userId })
        } catch (err) {
            this.error = err.message
            throw err
        }
    },
    async fetchMembers(groupId) {
        try {
            return await api.get(`/groups/${groupId}/members`)
        } catch (err) {
            console.error('Failed to fetch members:', err)
            throw err
        }
    },
    async fetchJoinRequests(groupId) {
        try {
            return await api.get(`/groups/${groupId}/requests`)
        } catch (err) {
            console.error('Failed to fetch join requests:', err)
            throw err
        }
    },
    async respondToRequest(requestId, accept) {
        try {
            await api.post(`/groups/requests/${requestId}/respond`, { accept })
        } catch (err) {
            console.error('Failed to respond to request:', err)
            throw err
        }
    },
    async fetchGroupEvents(groupId) {
        try {
            return await api.get(`/groups/${groupId}/events`)
        } catch (err) {
            console.error('Failed to fetch group events:', err)
            throw err
        }
    },
    async createEvent(groupId, eventData) {
        try {
            const dataToSubmit = { 
              ...eventData,
              event_time: new Date(eventData.event_time).toISOString()
            }
            return await api.post(`/groups/${groupId}/events`, dataToSubmit)
        } catch (err) {
            console.error('Failed to create event:', err)
            throw err
        }
    },
    async respondToEvent(eventId, response) {
        try {
            await api.post(`/groups/events/${eventId}/respond`, { response })
        } catch (err) {
            console.error('Failed to respond to event:', err)
            throw err
        }
    }
  }
})
