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
        
        // Determine follow status
        if (data.followers?.some(f => f.id === currentUserId)) {
          this.followStatus = 'accept'
        } else {
          this.followStatus = 'none'
        }
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
        const isAccepted = data.message.includes('accepted') || data.message.includes('following')
        this.followStatus = isAccepted ? 'accept' : 'pending'
        return isAccepted
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
