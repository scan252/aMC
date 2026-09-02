<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { NewsService } = bindings as any

interface RedeemCode {
  code: string
  reward: string
  source: string
  date: string
  expired: boolean
}

const codes = ref<RedeemCode[]>([])
const copied = ref('')

const REDEEM_URL = 'https://mc.appfeng.com/gachaLog' // 官方兑换入口（演示占位，待接入正式兑换页）

async function load() {
  codes.value = ((await NewsService.Codes()) as any[]) ?? []
}

function copy(code: string) {
  navigator.clipboard?.writeText(code)
  copied.value = code
  setTimeout(() => (copied.value = ''), 1500)
}

onMounted(load)
</script>

<template>
  <div class="page">
    <PageHeader title="兑换码" subtitle="新码速报 · 一键复制（数据源建设与社区渠道接入中）">
      <template #actions>
        <a class="btn" :href="REDEEM_URL" target="_blank">前往兑换</a>
      </template>
    </PageHeader>

    <div class="list">
      <div v-for="(c, i) in codes" :key="c.code" class="code glass" :class="{ expired: c.expired, [`a d${Math.min(i + 2, 6)}`]: true }">
        <div class="main">
          <b class="code-text">{{ c.code }}</b>
          <span class="reward">{{ c.reward }}</span>
        </div>
        <div class="meta">
          <span>{{ c.source }} · {{ c.date }}</span>
          <span v-if="c.expired" class="expired-tag">已过期</span>
        </div>
        <button class="btn" :disabled="c.expired" @click="copy(c.code)">
          {{ copied === c.code ? '已复制 ✓' : '复制' }}
        </button>
      </div>
      <div v-if="!codes.length" class="none glass">暂无兑换码</div>
    </div>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  overflow-y: auto;
}
.btn {
  height: 36px;
  padding: 0 18px;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  font-size: 12.5px;
  font-weight: 600;
  color: #201503;
  background: var(--gold-grad-btn);
  box-shadow: var(--gold-glow), inset 0 1px 0 rgba(255, 255, 255, 0.5);
  font-family: var(--font-ui);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}
.btn:disabled {
  background: rgba(255, 255, 255, 0.08);
  color: var(--txt-3);
  box-shadow: none;
  cursor: not-allowed;
}
.list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.code {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 16px 20px;
}
.code.expired {
  opacity: 0.55;
}
.main {
  flex: 1;
  min-width: 0;
}
.code-text {
  font-family: var(--font-mono);
  font-size: 16px;
  letter-spacing: 1px;
  display: block;
}
.reward {
  font-size: 12px;
  color: var(--warn);
}
.meta {
  font-size: 11.5px;
  color: var(--txt-3);
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}
.expired-tag {
  color: var(--danger);
}
.none {
  padding: 30px;
  text-align: center;
  color: var(--txt-3);
  font-size: 13px;
}
</style>
