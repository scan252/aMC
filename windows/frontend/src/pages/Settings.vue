<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="设置" subtitle="主题 · 数据目录 · 通知 · 自启" />
    <EmptyState icon="M4 8h9m4 0h3M4 16h3m4 0h9M15.5 5.6a2.4 2.4 0 1 1 0 4.8 2.4 2.4 0 0 1 0-4.8zM9 13.6a2.4 2.4 0 1 1 0 4.8 2.4 2.4 0 0 1 0-4.8z" title="设置" description="应用行为与数据管理。" />
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
}
</style>
