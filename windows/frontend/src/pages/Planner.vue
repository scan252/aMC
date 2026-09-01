<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="养成计算" subtitle="角色与武器养成材料规划" />
    <EmptyState icon="M5 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zM8 7h8M8 11h3m2 0h3M8 15h3m2 0h3" title="养成计算" description="选择目标等级，自动汇总突破与天赋材料清单，折算体力与预计天数。" />
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
