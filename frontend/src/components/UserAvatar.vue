<script setup>
import { computed } from 'vue'

const props = defineProps({
  url: String,
  name: String,
  size: {
    type: String,
    default: 'medium' // xsmall, small, medium, large
  }
})

const initials = computed(() => {
  return props.name ? props.name[0].toUpperCase() : 'U'
})
</script>

<template>
  <div class="avatar-wrapper" :class="size">
    <img v-if="url" :src="url" class="avatar-img" :alt="name" />
    <div v-else class="avatar-placeholder">{{ initials }}</div>
  </div>
</template>

<style scoped>
.avatar-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid var(--color-neon-cyan);
  background: var(--bg-accent);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-weight: bold;
  color: var(--color-neon-cyan);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Sizes */
.xsmall {
  width: 24px;
  height: 24px;
}
.xsmall .avatar-placeholder { font-size: 10px; }

.small {
  width: 40px;
  height: 40px;
}
.small .avatar-placeholder { font-size: 16px; }

.medium {
  width: 60px;
  height: 60px;
}
.medium .avatar-placeholder { font-size: 24px; }

.large {
  width: 120px;
  height: 120px;
}
.large .avatar-placeholder { font-size: 48px; }
</style>
