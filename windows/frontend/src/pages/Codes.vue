<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="兑换码" subtitle="新码速报 · 一键跳转兑换" />
    <EmptyState icon="M5 9h14a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V10a1 1 0 0 1 1-1zm7 1v11M4 9h16l-1.6-4.2a1 1 0 0 0-.94-.65H6.54a1 1 0 0 0-.94.65L4 9zm4.5-5.5C9.5 3 10 4 9.6 5.2L9.4 6M14.5 3.5C14.5 3 14 4 14.4 5.2l.2.8" title="兑换码" description="追踪官方与社区渠道的兑换码，新码发布即时提醒。" />
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
