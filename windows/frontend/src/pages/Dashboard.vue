<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '../components/PageHeader.vue'
import { useGachaStore } from '../stores/gacha'
import { useAccountStore } from '../stores/account'

const gacha = useGachaStore()
const account = useAccountStore()
const router = useRouter()

const hasGacha = computed(() => !!gacha.detail)
const stats = computed(() => gacha.detail?.stats ?? null)
const luckLabel = computed(() => {
  const l = stats.value?.luckIndex ?? 0
  if (!l) return '—'
  return l >= 1 ? `偏欧 ${Math.round((l - 1) * 100)}%` : `偏非 ${Math.round((1 - l) * 100)}%`
})
const wavePct = computed(() => {
  const w = account.waveplate
  if (!w?.max) return 0
  return Math.round((w.cur / w.max) * 100)
})

const quickLinks = [
  { to: '/gacha', name: '唤取分析', desc: '抽数 · 保底 · 分布', icon: 'M3 3v18h18M7 15l4-6 4 3 5-8' },
  { to: '/codex', name: '角色图鉴', desc: '资料库速查', icon: 'M4 19V5a2 2 0 0 1 2-2h13v16H6a2 2 0 0 0-2 2zm0 0a2 2 0 0 0 2 2h13' },
  { to: '/planner', name: '养成计算', desc: '材料清单规划', icon: 'M5 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zM8 7h8M8 11h3m2 0h3M8 15h3m2 0h3' },
]

onMounted(() => {
  gacha.init()
  account.load()
})
</script>

<template>
  <div class="page">
    <PageHeader title="总览" :subtitle="hasGacha ? `UID ${gacha.detail?.uid} · 更新于 ${gacha.detail?.fetchedAt}` : '配置你的账号数据，一切从这里开始'">
      <template #actions>
        <button class="btn" @click="router.push('/gacha')">进入唤取分析</button>
      </template>
    </PageHeader>

    <section v-if="hasGacha" class="stats">
      <div class="stat glass a d2">
        <div class="halo" />
        <div class="k">累计唤取</div>
        <div class="v gold-text">{{ stats?.total.toLocaleString() }}</div>
        <div class="foot">{{ stats?.pools.length ?? 0 }} 个卡池有记录</div>
      </div>
      <div class="stat glass a d2">
        <div class="k">5★ 已获取</div>
        <div class="v">{{ stats?.count5 }}<em>★</em></div>
        <div class="foot">4★ 共 {{ stats?.count4 }} 个</div>
      </div>
      <div class="stat glass a d3">
        <div class="k">平均出金</div>
        <div class="v">{{ stats?.avgPity ? stats.avgPity.toFixed(1) : '—' }}<em>抽</em></div>
        <div class="foot">综合期望 62.0 抽</div>
      </div>
      <div class="stat glass a d3">
        <div class="k">欧非指数</div>
        <div class="v">{{ stats?.luckIndex ? stats.luckIndex.toFixed(2) : '—' }}</div>
        <div class="foot">较均值<b>{{ luckLabel }}</b></div>
      </div>
    </section>

    <section class="mid">
      <!-- 账号中心快览 -->
      <div class="card glass a d3">
        <h3>账号快览<span class="tag">ACCOUNT</span></h3>
        <template v-if="account.bound">
          <div class="ov-grid">
            <div class="ov-item">
              <span class="k">签到状态</span>
              <b :class="account.overview?.signIn?.hadSignIn ? 'ok' : 'warn'">
                {{ account.overview?.signIn?.hadSignIn ? '已签到' : '未签到' }}
              </b>
            </div>
            <div class="ov-item">
              <span class="k">结晶波片</span>
              <b>{{ account.waveplate?.cur ?? '—' }}<small>/ {{ account.waveplate?.max ?? 240 }}</small></b>
            </div>
            <div class="ov-item">
              <span class="k">回满还需</span>
              <b>{{ account.waveFullInMin }}<small>分钟</small></b>
            </div>
          </div>
          <div class="wave-track">
            <div class="wave-fill" :style="{ width: `${wavePct}%` }" />
          </div>
        </template>
        <template v-else>
          <div class="bind-hint">
            <p>绑定库街区账号后，可使用每日自动签到、体力监控与练度查询。</p>
            <button class="btn ghost" @click="router.push('/account')">去绑定</button>
          </div>
        </template>
      </div>

      <!-- 唤取快览 -->
      <div v-if="hasGacha" class="card glass a d4">
        <h3>最近 5★<span class="tag">RECENT</span></h3>
        <div class="recent">
          <div v-for="r in gacha.recent5Filtered.slice(0, 5)" :key="r.time + r.name" class="recent-row">
            <span class="name">{{ r.name }}</span>
            <span class="meta">{{ r.time.slice(5, 10) }}<template v-if="r.gap >= 0"> · {{ r.gap }} 抽</template></span>
          </div>
        </div>
        <button class="link-more" @click="router.push('/gacha')">查看全部 →</button>
      </div>
      <div v-else class="card glass a d4">
        <h3>唤取数据<span class="tag">GACHA</span></h3>
        <div class="bind-hint">
          <p>还没有本地唤取数据。启动游戏打开唤取记录页后点击抓取，或直接导入 aMC JSON 数据文件（支持 Mac 版格式）。</p>
          <button class="btn ghost" @click="router.push('/gacha')">去抓取</button>
        </div>
      </div>

      <!-- 快捷入口 -->
      <div class="card glass a d5">
        <h3>快捷入口<span class="tag">TOOLS</span></h3>
        <div class="links">
          <div v-for="l in quickLinks" :key="l.to" class="link-row" @click="router.push(l.to)">
            <svg viewBox="0 0 24 24"><path :d="l.icon" /></svg>
            <div>
              <b>{{ l.name }}</b>
              <small>{{ l.desc }}</small>
            </div>
          </div>
        </div>
      </div>
    </section>
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
  transition: all 0.2s;
}
.btn:hover {
  transform: translateY(-1px);
}
.btn.ghost {
  background: rgba(255, 255, 255, 0.06);
  color: var(--txt-2);
  box-shadow: none;
  border: 1px solid var(--line);
}
.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  flex-shrink: 0;
}
.stat {
  padding: 20px 22px;
  position: relative;
  overflow: hidden;
}
.stat .halo {
  position: absolute;
  right: -34px;
  top: -34px;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--gold-soft), transparent 70%);
}
.stat .k {
  font-size: 12px;
  color: var(--txt-3);
  letter-spacing: 0.14em;
}
.stat .v {
  font-size: 38px;
  font-weight: 650;
  margin-top: 10px;
  letter-spacing: -1px;
}
.stat .v em {
  font-style: normal;
  font-size: 16px;
  font-weight: 500;
  color: var(--txt-3);
  margin-left: 4px;
}
.stat .foot {
  margin-top: 9px;
  font-size: 12px;
  color: var(--txt-2);
}
.stat .foot b {
  color: var(--ok);
}
.mid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 14px;
  flex: 1;
  min-height: 230px;
}
.card {
  padding: 20px 22px;
  display: flex;
  flex-direction: column;
}
.card h3 {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--txt-2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}
.card h3 .tag {
  font-size: 10.5px;
  color: var(--txt-3);
  letter-spacing: 0.08em;
}
.ov-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
  flex: 1;
}
.ov-item .k {
  display: block;
  font-size: 11.5px;
  color: var(--txt-3);
  margin-bottom: 5px;
}
.ov-item b {
  font-size: 19px;
  font-weight: 650;
}
.ov-item b small {
  font-size: 12px;
  font-weight: 500;
  color: var(--txt-3);
  margin-left: 2px;
}
.ov-item .ok {
  color: var(--ok);
}
.ov-item .warn {
  color: var(--warn);
}
.wave-track {
  height: 7px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.07);
  overflow: hidden;
  margin-top: 14px;
}
.wave-fill {
  height: 100%;
  border-radius: 5px;
  background: linear-gradient(90deg, var(--gold-3), var(--gold-1));
  box-shadow: 0 0 12px rgba(244, 206, 106, 0.5);
  transition: width 0.6s cubic-bezier(0.2, 0.7, 0.2, 1);
}
.bind-hint {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 14px;
}
.bind-hint p {
  font-size: 12.5px;
  color: var(--txt-3);
  line-height: 1.8;
}
.recent {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}
.recent-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
  border-bottom: 1px solid var(--line-2);
  font-size: 13px;
}
.recent-row:last-child {
  border-bottom: none;
}
.recent-row .name {
  font-weight: 600;
}
.recent-row .meta {
  color: var(--txt-3);
  font-size: 12px;
  font-family: var(--font-mono);
}
.link-more {
  background: none;
  border: none;
  color: var(--warn);
  font-size: 12.5px;
  cursor: pointer;
  text-align: left;
  padding: 10px 4px 0;
  font-family: var(--font-ui);
}
.links {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}
.link-row {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 11px 13px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--line-2);
  cursor: pointer;
  transition: all 0.2s;
}
.link-row:hover {
  background: rgba(255, 255, 255, 0.07);
  border-color: var(--line);
}
.link-row svg {
  width: 19px;
  height: 19px;
  stroke: var(--warn);
  fill: none;
  stroke-width: 1.6;
  stroke-linecap: round;
  stroke-linejoin: round;
  flex-shrink: 0;
}
.link-row b {
  display: block;
  font-size: 13.5px;
}
.link-row small {
  color: var(--txt-3);
  font-size: 11.5px;
}
</style>
