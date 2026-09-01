<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="角色图鉴" subtitle="角色 · 武器 · 声骸资料库" />
    <EmptyState icon="M4 19V5a2 2 0 0 1 2-2h13v16H6a2 2 0 0 0-2 2zm0 0a2 2 0 0 0 2 2h13" title="角色图鉴" description="内置鸣潮游戏数据库，支持按属性、武器类型与共鸣属性筛选。游戏数据管道建设中。" />
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
