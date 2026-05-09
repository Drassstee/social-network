import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: (() => {
      try {
        const u = localStorage.getItem('user')
        return u ? JSON.parse(u) : null
      } catch (e) {
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
        const response = await fetch('/api/v1/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password }),
        })

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}))
          throw new Error(errorData.error || 'Login failed')
        }

        const data = await response.json()
        
        this.user = data
        this.isAuthenticated = true
        localStorage.setItem('user', JSON.stringify(data))
        return true
      } catch (err) {
        this.error = err.message
        return false
      } finally {
        this.loading = false
      }
    },
    async register(userData) {
      this.loading = true
      this.error = null
      try {
        const isFormData = userData instanceof FormData
        const options = {
          method: 'POST',
          body: isFormData ? userData : JSON.stringify(userData),
        }
        
        if (!isFormData) {
          options.headers = { 'Content-Type': 'application/json' }
        }

        const response = await fetch('/api/v1/register', options)
        
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}))
          throw new Error(errorData.error || 'Registration failed')
        }

        const data = await response.json()
        
        this.user = data
        this.isAuthenticated = true
        localStorage.setItem('user', JSON.stringify(data))
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
        const res = await fetch('/api/v1/logout', { method: 'POST' })
        if (!res.ok) console.error('Logout request failed')
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
        const response = await fetch('/api/v1/me')
        if (response.ok) {
          const data = await response.json()
          this.user = data
          this.isAuthenticated = true
          localStorage.setItem('user', JSON.stringify(data))
        } else {
          this.user = null
          this.isAuthenticated = false
          localStorage.removeItem('user')
        }
      } catch (err) {
        console.error('Session check failed:', err)
        this.user = null
        this.isAuthenticated = false
        localStorage.removeItem('user')
      }
    }
  }
})
