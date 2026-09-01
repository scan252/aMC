<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="资讯日历" subtitle="官方公告 · 版本日历 · 活动倒计时" />
    <EmptyState icon="M4 7h16a1 1 0 0 1 1 1v11a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V8a1 1 0 0 1 1-1zm4-4v4m8-4v4M3 11h18" title="资讯日历" description="聚合官方公告与活动日程，卡池开始结束、活动截止一目了然。" />
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
