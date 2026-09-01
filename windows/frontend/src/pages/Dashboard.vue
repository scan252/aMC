<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="总览" subtitle="账号数据、待办事项与动态一屏尽览" />
    <EmptyState icon="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z" title="总览仪表盘" description="汇总你的唤取、签到、体力与版本动态。此模块将在账号域完成后接入实时数据。" />
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
