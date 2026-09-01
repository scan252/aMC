<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
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
</style>
