<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="账号中心" subtitle="库街区账号绑定 · 签到 · 结晶波片" />
    <EmptyState icon="M12 3a4 4 0 1 1 0 8 4 4 0 0 1 0-8ZM4 21c1.5-4 4.5-6 8-6s6.5 2 8 6" title="账号中心" description="绑定库街区账号后可使用每日自动签到、体力监控与练度查询。底座建设中，登录能力待实测。" />
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
