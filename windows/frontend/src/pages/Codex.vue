<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { GameDataService } = bindings as any

interface Character {
  id: string
  name: string
  rarity: number
  element: string
  weapon: string
}

const keyword = ref('')
const element = ref('')
const weapon = ref('')
const chars = ref<Character[]>([])

const elements = ['衍射', '湮灭', '气动', '导电', '冷凝', '热熔']
const weapons = ['长刃', '迅刀', '佩枪', '臂铠', '音感仪']

const elementColor: Record<string, string> = {
  衍射: '#F6C87E', 湮灭: '#BCA8F7', 气动: '#7DE0C3',
  导电: '#A8C8F7', 冷凝: '#9AD5F0', 热熔: '#F5907E',
}

async function load() {
  chars.value = ((await GameDataService.Characters(keyword.value, element.value, weapon.value)) as any[]) ?? []
}

onMounted(load)
</script>

<template>
  <div class="page">
    <PageHeader title="角色图鉴" subtitle="内置数据库种子版 · 随版本更新通道扩充" />

    <div class="filters glass a d1">
      <input v-model="keyword" class="field" placeholder="搜索角色名…" @input="load" />
      <select v-model="element" class="field" @change="load">
        <option value="">全部属性</option>
        <option v-for="e in elements" :key="e" :value="e">{{ e }}</option>
      </select>
      <select v-model="weapon" class="field" @change="load">
        <option value="">全部武器</option>
        <option v-for="w in weapons" :key="w" :value="w">{{ w }}</option>
      </select>
    </div>

    <div class="grid">
      <div v-for="c in chars" :key="c.id" class="char glass a d3">
        <div class="avatar" :style="{ background: `radial-gradient(circle at 30% 20%, ${elementColor[c.element] ?? '#888'}33, transparent 70%)` }">
          <b>{{ c.name.slice(0, 1) }}</b>
        </div>
        <div class="info">
          <div class="name-row">
            <b>{{ c.name }}</b>
            <span class="rarity">{{ '★'.repeat(c.rarity) }}</span>
          </div>
          <div class="tags">
            <span class="chip" :style="{ color: elementColor[c.element], borderColor: `${elementColor[c.element]}55` }">{{ c.element }}</span>
            <span class="chip dim">{{ c.weapon }}</span>
          </div>
        </div>
      </div>
    </div>
    <div v-if="!chars.length" class="none glass">没有符合条件的角色</div>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
  overflow-y: auto;
}
.filters {
  display: flex;
  gap: 10px;
  padding: 12px 14px;
  flex-shrink: 0;
}
.field {
  height: 36px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-1);
  padding: 0 16px;
  font-size: 13px;
  outline: none;
  font-family: var(--font-ui);
}
.field:first-child {
  flex: 1;
}
.field:focus {
  border-color: rgba(237, 190, 90, 0.4);
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}
.char {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 14px 16px;
}
.avatar {
  width: 52px;
  height: 52px;
  border-radius: 15px;
  display: grid;
  place-items: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--line);
  flex-shrink: 0;
}
.avatar b {
  font-size: 20px;
  color: var(--txt-1);
}
.info {
  min-width: 0;
}
.name-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}
.name-row b {
  font-size: 15px;
}
.rarity {
  font-size: 10px;
  color: var(--gold-1);
  letter-spacing: 1px;
}
.tags {
  display: flex;
  gap: 6px;
  margin-top: 7px;
}
.chip {
  height: 22px;
  padding: 0 10px;
  border-radius: 999px;
  font-size: 11px;
  display: inline-flex;
  align-items: center;
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.04);
  color: var(--txt-2);
}
.none {
  padding: 30px;
  text-align: center;
  color: var(--txt-3);
  font-size: 13px;
}
</style>
