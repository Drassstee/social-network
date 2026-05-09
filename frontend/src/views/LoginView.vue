<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const email = ref('')
const password = ref('')
const auth = useAuthStore()
const router = useRouter()

const handleLogin = async () => {
  const success = await auth.login(email.value, password.value)
  if (success) {
    router.push('/')
  }
}
</script>

<template>
  <div class="login-container">
    <div class="card-retro login-card">
      <div class="card-header">
        <h2 class="title">AUTH_SYSTEM</h2>
        <p class="subtitle">INPUT_CREDENTIALS</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email">USER_IDENTIFIER</label>
          <input 
            id="email"
            v-model="email" 
            type="email" 
            class="input-retro" 
            placeholder="USER@NET.SYS"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">ACCESS_CODE</label>
          <input 
            id="password"
            v-model="password" 
            type="password" 
            class="input-retro" 
            placeholder="********"
            required
          />
        </div>
        
        <div v-if="auth.error" class="error-msg">
          ERROR: {{ auth.error }}
        </div>
        
        <button 
          type="submit" 
          class="btn-retro w-full"
          :disabled="auth.loading"
        >
          {{ auth.loading ? 'INITIALIZING...' : 'EXECUTE_LOGIN' }}
        </button>
      </form>
      
      <div class="card-footer">
        <span>NEW_USER? </span>
        <router-link to="/register" class="link-retro">REGISTER_ID</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 80px);
}

.login-card {
  width: 100%;
  max-width: 450px;
  animation: glitch-in 0.5s ease-out;
}

@keyframes glitch-in {
  0% { transform: scale(0.9); opacity: 0; filter: hue-rotate(90deg); }
  100% { transform: scale(1); opacity: 1; filter: hue-rotate(0deg); }
}

.card-header {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  font-size: 2.2rem;
  margin-bottom: 5px;
  color: var(--color-neon-cyan);
  text-shadow: var(--shadow-neon-cyan);
}

.subtitle {
  color: var(--color-neon-magenta);
  font-family: 'VT323', monospace;
  font-size: 1.5rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

label {
  font-weight: 700;
  color: var(--color-neon-cyan);
  font-size: 1.1rem;
  font-family: 'VT323', monospace;
}

.w-full {
  width: 100%;
  padding: 15px;
  font-size: 1.2rem;
}

.error-msg {
  color: var(--color-neon-yellow);
  background: rgba(255, 255, 0, 0.1);
  padding: 10px;
  border: 1px solid var(--color-neon-yellow);
  font-size: 1.1rem;
  text-align: center;
  font-family: 'VT323', monospace;
}

.card-footer {
  margin-top: 30px;
  text-align: center;
  font-size: 1.1rem;
  font-family: 'VT323', monospace;
}

.link-retro {
  color: var(--color-neon-magenta);
  text-decoration: none;
  font-weight: 700;
}

.link-retro:hover {
  text-shadow: 0 0 5px var(--color-neon-magenta);
  text-decoration: underline;
}
</style>
