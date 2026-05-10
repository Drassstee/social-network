<script setup>
import { useUIStore } from '../stores/ui'

const ui = useUIStore()
</script>

<template>
  <Transition name="toast">
    <div v-if="ui.toast.visible" class="feedback-toast" :class="ui.toast.type" @click="ui.hideToast">
      {{ ui.toast.message }}
    </div>
  </Transition>
</template>

<style scoped>
.feedback-toast {
  position: fixed;
  bottom: 30px;
  right: 30px;
  padding: 15px 25px;
  background: var(--color-dark-bg);
  border: 2px solid var(--color-neon-cyan);
  color: var(--color-neon-cyan);
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
  z-index: 1000;
  box-shadow: 5px 5px 0 var(--color-neon-magenta);
  cursor: pointer;
  max-width: 300px;
}

.feedback-toast.success {
  border-color: #00ff00;
  color: #00ff00;
  box-shadow: 5px 5px 0 rgba(0, 255, 0, 0.3);
}

.feedback-toast.error {
  border-color: #ff0055;
  color: #ff0055;
  box-shadow: 5px 5px 0 rgba(255, 0, 85, 0.3);
}

/* Vue Transitions */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  transform: translateY(30px);
  opacity: 0;
}

.toast-leave-to {
  transform: translateX(100px);
  opacity: 0;
}
</style>
