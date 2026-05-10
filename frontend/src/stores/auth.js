import { defineStore } from 'pinia'
import { api } from '../api/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: (() => {
      try {
        const u = localStorage.getItem('user')
        return u ? JSON.parse(u) : null
      } catch {
        localStorage.removeItem('user')
        return null
      }
    })(),
    isAuthenticated: !!localStorage.getItem('user'),
    loading: false,
    error: null,
  }),
  actions: {
    async login(email, password) {
      this.loading = true
      this.error = null
      try {
        const data = await api.post('/login', { email, password })
        this.setUser(data)
        return true
      } catch (err) {
        this.error = err.message
        return false
      } finally {
        this.loading = false
      }
    },
    async register(payload, avatarFile = null) {
      this.loading = true
      this.error = null
      try {
        const dataToSubmit = { ...payload }
        if (dataToSubmit.date_of_birth) {
          dataToSubmit.dob = new Date(dataToSubmit.date_of_birth).toISOString()
          delete dataToSubmit.date_of_birth
        }

        let body = dataToSubmit
        if (avatarFile) {
          body = new FormData()
          for (const key in dataToSubmit) {
            body.append(key, dataToSubmit[key])
          }
          body.append('avatar', avatarFile)
        }

        const data = await api.post('/register', body)
        this.setUser(data)
        return true
      } catch (err) {
        this.error = err.message
        return false
      } finally {
        this.loading = false
      }
    },
    async logout() {
      try {
        await api.post('/logout')
      } catch (err) {
        console.error('Logout error:', err)
      } finally {
        const { useChatStore } = await import('./chat')
        const chatStore = useChatStore()
        chatStore.disconnect()
        
        this.user = null
        this.isAuthenticated = false
        localStorage.removeItem('user')
      }
    },
    async checkSession() {
      try {
        const data = await api.get('/me')
        this.setUser(data)
      } catch (err) {
        if (err.status !== 401) {
          console.error('Session check failed:', err)
        }
        this.clearUser()
      }
    },
    setUser(data) {
      this.user = data
      this.isAuthenticated = true
      localStorage.setItem('user', JSON.stringify(data))
    },
    clearUser() {
      this.user = null
      this.isAuthenticated = false
      localStorage.removeItem('user')
    }
  }
})
