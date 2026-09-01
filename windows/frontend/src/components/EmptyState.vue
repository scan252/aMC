<script setup lang="ts">
defineProps<{
  icon: string
  title: string
  description: string
}>()
</script>

<template>
  <div class="empty glass a d3">
    <div class="icon-wrap">
      <svg viewBox="0 0 24 24"><path :d="icon" /></svg>
    </div>
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
    <slot />
  </div>
</template>

<style scoped>
.empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px;
  min-height: 0;
}
.icon-wrap {
  width: 74px;
  height: 74px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: var(--gold-soft);
  border: 1px solid rgba(237, 190, 90, 0.2);
  margin-bottom: 22px;
}
.icon-wrap svg {
  width: 32px;
  height: 32px;
  stroke: var(--warn);
  fill: none;
  stroke-width: 1.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}
h2 {
  font-size: 17px;
  font-weight: 650;
  margin-bottom: 9px;
}
p {
  font-size: 13px;
  color: var(--txt-3);
  max-width: 400px;
  line-height: 1.7;
}
</style>
