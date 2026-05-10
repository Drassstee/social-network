import { defineStore } from 'pinia'

export const useUIStore = defineStore('ui', {
  state: () => ({
    toast: {
      message: '',
      type: 'info', // 'info', 'success', 'error'
      visible: false
    }
  }),
  actions: {
    showToast(message, type = 'info') {
      this.toast.message = message
      this.toast.type = type
      this.toast.visible = true

      setTimeout(() => {
        this.toast.visible = false
      }, 3000)
    },
    hideToast() {
      this.toast.visible = false
    }
  }
})
