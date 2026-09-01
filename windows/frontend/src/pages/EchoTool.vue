<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="声骸工具" subtitle="声骸评分 · 毕业度分析" />
    <EmptyState icon="M12 3l8 4.5v9L12 21l-8-4.5v-9L12 3zm0 9l8-4.5M12 12v9M12 12L4 7.5" title="声骸工具" description="导入或录入声骸词条，自动计算评分与毕业度，找到最值得强化的那一件。" />
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
