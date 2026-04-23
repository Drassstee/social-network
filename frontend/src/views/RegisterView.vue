<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const userData = ref({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
  date_of_birth: '',
  nickname: '',
  about_me: '',
  avatar_url: ''
})

const auth = useAuthStore()
const router = useRouter()

const handleRegister = async () => {
  const payload = { ...userData.value }
  if (payload.date_of_birth) {
    payload.dob = new Date(payload.date_of_birth).toISOString()
    delete payload.date_of_birth
  }
  const success = await auth.register(payload)
  if (success) {
    router.push('/')
  }
}
</script>

<template>
  <div class="register-container">
    <div class="card-traditional register-card">
      <div class="card-header">
        <h2 class="title">Join Our Circle</h2>
        <p class="subtitle">New Member Registration</p>
      </div>
      
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-row">
          <div class="form-group">
            <label for="first_name">First Name</label>
            <input 
              id="first_name"
              v-model="userData.first_name" 
              type="text" 
              class="input-traditional" 
              required
            />
          </div>
          <div class="form-group">
            <label for="last_name">Last Name</label>
            <input 
              id="last_name"
              v-model="userData.last_name" 
              type="text" 
              class="input-traditional" 
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="email">Email</label>
          <input 
            id="email"
            v-model="userData.email" 
            type="email" 
            class="input-traditional" 
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            id="password"
            v-model="userData.password" 
            type="password" 
            class="input-traditional" 
            required
          />
        </div>

        <div class="form-group">
          <label for="dob">Date of Birth</label>
          <input 
            id="dob"
            v-model="userData.date_of_birth" 
            type="date" 
            class="input-traditional" 
            required
          />
        </div>

        <div class="form-divider">
          <span>Optional Information</span>
        </div>

        <div class="form-group">
          <label for="nickname">Nickname (Optional)</label>
          <input 
            id="nickname"
            v-model="userData.nickname" 
            type="text" 
            class="input-traditional" 
          />
        </div>

        <div class="form-group">
          <label for="about">About Me (Optional)</label>
          <textarea 
            id="about"
            v-model="userData.about_me" 
            class="input-traditional textarea" 
            rows="3"
          ></textarea>
        </div>
        
        <div v-if="auth.error" class="error-msg">
          {{ auth.error }}
        </div>
        
        <button 
          type="submit" 
          class="btn-traditional w-full"
          :disabled="auth.loading"
        >
          {{ auth.loading ? 'Creating account...' : 'Register' }}
        </button>
      </form>
      
      <div class="card-footer">
        <span>Already have an account? </span>
        <router-link to="/login" class="link-vermilion">Sign in</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
}

.register-card {
  width: 100%;
  max-width: 600px;
  animation: fadeInUp 0.8s ease-out;
}

.card-header {
  text-align: center;
  margin-bottom: 30px;
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

.register-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin: 10px 0;
  color: #888;
  font-size: 0.85rem;
}

.form-divider::before, .form-divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid #ddd;
}

.form-divider:not(:empty)::before {
  margin-right: 15px;
}

.form-divider:not(:empty)::after {
  margin-left: 15px;
}

label {
  font-weight: 500;
  color: var(--color-charcoal);
  font-size: 0.9rem;
}

.textarea {
  resize: vertical;
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

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 600px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
