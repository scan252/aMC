<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { Events } from '@wailsio/runtime'
import Sidebar from './components/Sidebar.vue'

interface NotifyPayload {
  title: string
  body: string
}

const toasts = ref<(NotifyPayload & { id: number })[]>([])
let toastId = 0
let cancelNotify: (() => void) | null = null

function pushToast(p: NotifyPayload) {
  const id = ++toastId
  toasts.value.push({ ...p, id })
  setTimeout(() => {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }, 5200)
}

onMounted(() => {
  cancelNotify = Events.On('app:notify', (ev: any) => {
    const data = ev?.data ?? ev
    if (data?.title) pushToast(data as NotifyPayload)
  }) as unknown as () => void
})

onBeforeUnmount(() => {
  cancelNotify?.()
})
</script>

<template>
  <div class="orb orb-a" />
  <div class="orb orb-b" />
  <div class="orb orb-c" />
  <div class="shell">
    <Sidebar />
    <main class="content">
      <router-view v-slot="{ Component }">
        <transition name="fade-slide" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>

  <!-- 应用内通知 -->
  <div class="toasts">
    <transition-group name="toast">
      <div v-for="t in toasts" :key="t.id" class="toast glass">
        <b>{{ t.title }}</b>
        <p>{{ t.body }}</p>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.shell {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: 236px 1fr;
  gap: 16px;
  height: 100vh;
  padding: 16px;
}
.content {
  min-width: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.toasts {
  position: fixed;
  right: 20px;
  top: 20px;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 300px;
}
.toast {
  padding: 13px 17px;
  border-radius: var(--radius-sm);
}
.toast b {
  display: block;
  font-size: 13px;
  margin-bottom: 4px;
  color: var(--warn);
}
.toast p {
  font-size: 12px;
  color: var(--txt-2);
  line-height: 1.6;
}
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
