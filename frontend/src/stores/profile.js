import { defineStore } from 'pinia'
import { api } from '../api/api'

export const useProfileStore = defineStore('profile', {
  state: () => ({
    profile: null,
    loading: false,
    error: null,
    followStatus: 'none',
    myFollowers: []
  }),
  actions: {
    async fetchMyFollowers(userId) {
      if (!userId) return
      try {
        const data = await api.get(`/users/${userId}`)
        this.myFollowers = data.followers || []
        return this.myFollowers
      } catch (err) {
        console.error('Failed to fetch followers:', err)
        throw err
      }
    },
    async fetchProfile(id, currentUserId) {
      this.loading = true
      try {
        const data = await api.get(`/users/${id}`)
        this.profile = data
        
        // Use the follow_status provided by the backend
        this.followStatus = data.follow_status || 'none'
        
        return data
      } catch (err) {
        this.error = err.message
        this.profile = null
        throw err
      } finally {
        this.loading = false
      }
    },
    async follow(targetId) {
      try {
        const data = await api.post('/follow', { following_id: Number(targetId) })
        // Backend returns "accept" or "pending" as message
        this.followStatus = data.message === 'accept' ? 'following' : 'pending'
        return this.followStatus
      } catch (err) {
        console.error('Follow failed:', err)
        throw err
      }
    },
    async unfollow(targetId) {
      try {
        await api.post('/unfollow', { following_id: Number(targetId) })
        this.followStatus = 'none'
      } catch (err) {
        console.error('Unfollow failed:', err)
        throw err
      }
    },
    async updateProfile(userData) {
      try {
        const data = await api.put('/users', userData)
        if (this.profile && this.profile.user.id === userData.id) {
          this.profile.user = { ...this.profile.user, ...userData }
        }
        return data
      } catch (err) {
        console.error('Profile update failed:', err)
        throw err
      }
    }
  }
})
