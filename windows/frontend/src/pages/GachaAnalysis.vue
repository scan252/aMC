<script setup lang="ts">
import PageHeader from '../components/PageHeader.vue'
</script>

<template>
  <div class="page">
    <PageHeader title="唤取分析" subtitle="数据更新于 — · 来自官方唤取接口">
      <template #actions>
        <div class="account glass">UID — · 国服 ▾</div>
        <button class="btn">立即抓取</button>
      </template>
    </PageHeader>
    <div class="body glass a d2">
      <div class="center">
        <div class="icon-wrap">
          <svg viewBox="0 0 24 24">
            <path d="M3 3v18h18M7 15l4-6 4 3 5-8" />
          </svg>
        </div>
        <h2>暂无唤取数据</h2>
        <p>
          启动游戏并打开「唤取 → 唤取记录」页面后，点击「立即抓取」。
          <br />
          抓取引擎正在接入中——日志发现、解密与官方接口拉取将在后续提交完成。
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
}
.body {
  flex: 1;
  display: flex;
  min-height: 0;
}
.center {
  margin: auto;
  text-align: center;
  padding: 40px;
}
.icon-wrap {
  width: 74px;
  height: 74px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: var(--gold-soft);
  border: 1px solid rgba(237, 190, 90, 0.2);
  margin: 0 auto 22px;
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
  line-height: 1.8;
}
.account {
  display: flex;
  align-items: center;
  gap: 9px;
  height: 38px;
  padding: 0 15px;
  border-radius: 999px;
  font-size: 12.5px;
  color: var(--txt-2);
  cursor: pointer;
}
.btn {
  height: 38px;
  padding: 0 20px;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  color: #201503;
  background: var(--gold-grad-btn);
  letter-spacing: 0.02em;
  box-shadow: var(--gold-glow), inset 0 1px 0 rgba(255, 255, 255, 0.5);
  transition: all 0.22s;
}
.btn:hover {
  transform: translateY(-1px);
}
</style>
