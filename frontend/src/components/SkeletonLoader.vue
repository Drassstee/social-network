<script setup>
defineProps({
  type: {
    type: String,
    default: 'card' // 'card', 'text', 'avatar', 'list-item'
  },
  count: {
    type: Number,
    default: 1
  }
})
</script>

<template>
  <div class="skeleton-container">
    <div v-for="i in count" :key="i" :class="['skeleton-item', type]">
      <div v-if="type === 'card'" class="card-skeleton">
        <div class="header-skel">
          <div class="avatar-skel pulse"></div>
          <div class="meta-skel">
            <div class="title-skel pulse"></div>
            <div class="subtitle-skel pulse"></div>
          </div>
        </div>
        <div class="content-skel pulse"></div>
        <div class="footer-skel pulse"></div>
      </div>
      
      <div v-else-if="type === 'avatar'" class="avatar-skel pulse big"></div>
      
      <div v-else-if="type === 'list-item'" class="list-item-skel">
        <div class="icon-skel pulse"></div>
        <div class="text-skel pulse"></div>
      </div>
      
      <div v-else class="text-skel pulse"></div>
    </div>
  </div>
</template>

<style scoped>
.skeleton-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.skeleton-item {
  width: 100%;
}

.pulse {
  background: linear-gradient(
    90deg,
    rgba(0, 255, 255, 0.05) 25%,
    rgba(0, 255, 255, 0.1) 50%,
    rgba(0, 255, 255, 0.05) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border: 1px solid rgba(0, 255, 255, 0.1);
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Card Skeleton */
.card-skeleton {
  padding: 20px;
  border: 2px solid rgba(0, 255, 255, 0.1);
  background: rgba(31, 11, 53, 0.5);
  box-shadow: 4px 4px 0 rgba(255, 0, 255, 0.1);
}

.header-skel {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.avatar-skel {
  width: 45px;
  height: 45px;
  flex-shrink: 0;
}

.avatar-skel.big {
  width: 120px;
  height: 120px;
  margin: 0 auto;
}

.meta-skel {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.title-skel {
  height: 16px;
  width: 60%;
}

.subtitle-skel {
  height: 12px;
  width: 30%;
}

.content-skel {
  height: 80px;
  width: 100%;
  margin-bottom: 15px;
}

.footer-skel {
  height: 20px;
  width: 40%;
}

/* List Item Skeleton */
.list-item-skel {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border: 1px solid rgba(0, 255, 255, 0.1);
}

.icon-skel {
  width: 40px;
  height: 40px;
}

.text-skel {
  height: 14px;
  width: 70%;
}
</style>
