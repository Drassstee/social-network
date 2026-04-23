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
    <div class="card-traditional login-card">
      <div class="card-header">
        <h2 class="title">Welcome Back</h2>
        <p class="subtitle">Welcome Back</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            id="email"
            v-model="email" 
            type="email" 
            class="input-traditional" 
            placeholder="example@mail.com"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            id="password"
            v-model="password" 
            type="password" 
            class="input-traditional" 
            placeholder="••••••••"
            required
          />
        </div>
        
        <div v-if="auth.error" class="error-msg">
          {{ auth.error }}
        </div>
        
        <button 
          type="submit" 
          class="btn-traditional w-full"
          :disabled="auth.loading"
        >
          {{ auth.loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>
      
      <div class="card-footer">
        <span>Don't have an account? </span>
        <router-link to="/register" class="link-vermilion">Register now</router-link>
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
  animation: fadeInDown 0.8s ease-out;
}

.card-header {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  font-size: 2rem;
  margin-bottom: 5px;
}

.subtitle {
  color: var(--color-vermilion);
  font-family: 'Noto Serif JP', serif;
  font-size: 1.2rem;
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
  font-weight: 500;
  color: var(--color-charcoal);
  font-size: 0.9rem;
}

.w-full {
  width: 100%;
  padding: 15px;
  font-size: 1.1rem;
}

.error-msg {
  color: var(--color-vermilion);
  background: rgba(188, 0, 45, 0.1);
  padding: 10px;
  border-radius: var(--border-radius);
  font-size: 0.85rem;
  text-align: center;
}

.card-footer {
  margin-top: 30px;
  text-align: center;
  font-size: 0.9rem;
}

.link-vermilion {
  color: var(--color-vermilion);
  text-decoration: none;
  font-weight: 600;
}

.link-vermilion:hover {
  text-decoration: underline;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
