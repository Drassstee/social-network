<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="not-found-view">
    <div class="card-retro error-card">
      <div class="glitch-wrapper">
        <h1 class="error-code" data-text="404">404</h1>
      </div>
      <h2 class="error-title">SEGMENTATION_FAULT</h2>
      <p class="error-msg">
        THE REQUESTED MEMORY ADDRESS (URL) IS UNREACHABLE OR HAS BEEN WIPED FROM THE MAINFRAME.
      </p>
      
      <div class="visual-deco">
        <div class="static-line"></div>
        <div class="static-line"></div>
        <div class="static-line"></div>
      </div>

      <button @click="goHome" class="btn-retro large">
        REBOOT_TO_HOME
      </button>
    </div>
  </div>
</template>

<style scoped>
.not-found-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
  padding: 20px;
}

.error-card {
  text-align: center;
  max-width: 500px;
  padding: 60px 40px;
  position: relative;
  overflow: hidden;
}

.glitch-wrapper {
  position: relative;
}

.error-code {
  font-size: 8rem;
  font-weight: 900;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
  margin-bottom: 10px;
  font-family: 'Press Start 2P', cursive;
}

.error-title {
  color: var(--color-neon-magenta);
  font-size: 1.8rem;
  margin-bottom: 20px;
  letter-spacing: 2px;
}

.error-msg {
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
  font-size: 1.4rem;
  line-height: 1.4;
  margin-bottom: 40px;
}

.visual-deco {
  margin-bottom: 40px;
}

.static-line {
  height: 2px;
  background: var(--color-neon-cyan);
  margin: 10px auto;
  opacity: 0.3;
}

.static-line:nth-child(2) {
  width: 80%;
  background: var(--color-neon-magenta);
}

.static-line:nth-child(3) {
  width: 60%;
}

.btn-retro.large {
  width: 100%;
  font-size: 1.4rem;
  padding: 15px;
}

/* Optional Glitch Effect */
@keyframes glitch {
  0% { transform: translate(0); }
  20% { transform: translate(-2px, 2px); }
  40% { transform: translate(-2px, -2px); }
  60% { transform: translate(2px, 2px); }
  80% { transform: translate(2px, -2px); }
  100% { transform: translate(0); }
}

.error-code:hover {
  animation: glitch 0.3s infinite;
}
</style>
