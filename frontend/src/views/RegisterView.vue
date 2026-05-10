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
  about_me: ''
})

const auth = useAuthStore()
const router = useRouter()

const avatarFile = ref(null)

const handleAvatarChange = (e) => {
  avatarFile.value = e.target.files[0]
}

const handleRegister = async () => {
  const success = await auth.register(userData.value, avatarFile.value)
  if (success) {
    router.push('/')
  }
}
</script>

<template>
  <div class="register-container">
    <div class="card-retro register-card">
      <div class="card-header">
        <h2 class="title">NEW_IDENTITY_WIZARD</h2>
        <p class="subtitle">ALLOCATE_USER_RESOURCES</p>
      </div>
      
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-row">
          <div class="form-group">
            <label for="first_name">GIVEN_NAME</label>
            <input 
              id="first_name"
              v-model="userData.first_name" 
              type="text" 
              class="input-retro" 
              required
            />
          </div>
          <div class="form-group">
            <label for="last_name">SURNAME</label>
            <input 
              id="last_name"
              v-model="userData.last_name" 
              type="text" 
              class="input-retro" 
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="email">NET_ADDRESS</label>
          <input 
            id="email"
            v-model="userData.email" 
            type="email" 
            class="input-retro" 
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">ENCRYPTION_KEY</label>
          <input 
            id="password"
            v-model="userData.password" 
            type="password" 
            class="input-retro" 
            required
          />
        </div>

        <div class="form-group">
          <label for="dob">TIMESTAMP_OF_ORIGIN</label>
          <input 
            id="dob"
            v-model="userData.date_of_birth" 
            type="date" 
            class="input-retro" 
            required
          />
        </div>

        <div class="form-divider">
          <span>EXTENDED_PARAMETERS</span>
        </div>

        <div class="form-group">
          <label for="nickname">ALIAS (OPTIONAL)</label>
          <input 
            id="nickname"
            v-model="userData.nickname" 
            type="text" 
            class="input-retro" 
          />
        </div>

        <div class="form-group">
          <label for="avatar">VISUAL_UID (OPTIONAL)</label>
          <input 
            id="avatar"
            type="file" 
            @change="handleAvatarChange" 
            accept="image/*"
            class="input-retro file-input" 
          />
        </div>

        <div class="form-group">
          <label for="about">BIO_DATA (OPTIONAL)</label>
          <textarea 
            id="about"
            v-model="userData.about_me" 
            class="input-retro textarea" 
            rows="3"
          ></textarea>
        </div>
        
        <div v-if="auth.error" class="error-msg">
          CRITICAL_ERROR: {{ auth.error }}
        </div>
        
        <button 
          type="submit" 
          class="btn-retro w-full"
          :disabled="auth.loading"
        >
          {{ auth.loading ? 'PROCESSING...' : 'ESTABLISH_ID' }}
        </button>
      </form>
      
      <div class="card-footer">
        <span>ID_EXISTS? </span>
        <router-link to="/login" class="link-retro">INITIATE_LOGIN</router-link>
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
  animation: slide-in 0.6s cubic-bezier(0.23, 1, 0.32, 1);
}

@keyframes slide-in {
  from { transform: translateX(-50px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

.card-header {
  text-align: center;
  margin-bottom: 30px;
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
  color: var(--color-neon-yellow);
  font-family: 'VT323', monospace;
  font-size: 1.2rem;
}

.form-divider::before, .form-divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--color-grid-line);
}

.form-divider:not(:empty)::before {
  margin-right: 15px;
}

.form-divider:not(:empty)::after {
  margin-left: 15px;
}

label {
  font-weight: 700;
  color: var(--color-neon-cyan);
  font-size: 1rem;
  font-family: 'VT323', monospace;
}

.textarea {
  resize: vertical;
}

.file-input {
  font-size: 1rem;
  padding: 8px;
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

@media (max-width: 600px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
