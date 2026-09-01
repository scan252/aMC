<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <header class="page-head glass a d1">
    <div>
      <h1>{{ title }}</h1>
      <p v-if="subtitle" class="sub">{{ subtitle }}</p>
    </div>
    <div class="actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  flex-shrink: 0;
}
h1 {
  font-size: 22px;
  font-weight: 650;
  letter-spacing: 0.01em;
}
.sub {
  font-size: 12.5px;
  color: var(--txt-3);
  margin-top: 5px;
}
.actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
