<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import EmptyState from '../components/EmptyState.vue'
import { useGachaStore, POOL_NAMES } from '../stores/gacha'

const gacha = useGachaStore()
const accountMenuOpen = ref(false)

const activeDetail = computed(() => gacha.detail)
const activeStats = computed(() => activeDetail.value?.stats ?? null)
const activePool = computed({
  get: () => gacha.activePool,
  set: (v: string) => (gacha.activePool = v),
})
const poolTabs = computed(() => {
  const tabs: { key: string; name: string; count: number }[] = [
    { key: 'all', name: '全部卡池', count: activeStats.value?.total ?? 0 },
  ]
  for (const p of activeStats.value?.pools ?? []) {
    tabs.push({ key: String(p.poolType), name: p.poolName, count: p.total })
  }
  return tabs
})
const activePoolName = computed(() =>
  gacha.activePool === 'all' ? '全部卡池' : POOL_NAMES[gacha.activePool] ?? gacha.activePool,
)
const distMax = computed(() => Math.max(1, ...gacha.distFiltered))
const star5Split = computed(() => gacha.star5Split)
const star5Total = computed(() => {
  const s = gacha.star5Split
  return s.limited + s.standard
})
const limitedPct = computed(() => {
  const s = gacha.star5Split
  if (star5Total.value === 0) return 0
  return Math.round((s.limited / star5Total.value) * 100)
})
const luckLabel = computed(() => {
  const l = activeStats.value?.luckIndex ?? 0
  if (!l) return '—'
  return l >= 1 ? `偏欧 ${Math.round((l - 1) * 100)}%` : `偏非 ${Math.round((1 - l) * 100)}%`
})

const currentUid = computed(() => activeDetail.value?.uid ?? '')
const currentSvr = computed(() => (activeDetail.value?.svrArea === 'global' ? '国际服' : '国服'))
const updatedText = computed(() =>
  activeDetail.value?.fetchedAt ? `数据更新于 ${activeDetail.value.fetchedAt} · 来自官方唤取接口` : '尚无本地数据',
)

function fmtDate(time: string) {
  return time.slice(5, 10)
}
function distLabel(i: number) {
  return ['≤10', '11-20', '21-30', '31-40', '41-50', '51-60', '61-70', '71-75', '76-80'][i]
}
function poolNameOf(key: string) {
  return POOL_NAMES[key] ?? key
}

onMounted(() => gacha.init())
</script>

<template>
  <div class="page">
    <PageHeader title="唤取分析" :subtitle="updatedText">
      <template #actions>
        <div v-if="gacha.accounts.length" class="account glass" @click="accountMenuOpen = !accountMenuOpen">
          <b>UID {{ currentUid || gacha.accounts[0].uid }}</b>{{ currentSvr }}
          <span class="caret" />
          <div v-if="accountMenuOpen" class="menu glass">
            <div
              v-for="acc in gacha.accounts"
              :key="acc.uid"
              class="menu-item"
              :class="{ on: acc.uid === currentUid }"
              @click.stop="gacha.switchAccount(acc.uid); accountMenuOpen = false"
            >
              <span>UID {{ acc.uid }}</span>
              <small>{{ acc.total }} 抽 · {{ acc.count5 }}★</small>
            </div>
          </div>
        </div>
        <button class="btn" :disabled="gacha.loading" @click="gacha.refresh()">
          {{ gacha.loading ? '抓取中…' : '立即抓取' }}
        </button>
      </template>
    </PageHeader>

    <!-- 抓取进度条 -->
    <div v-if="gacha.loading" class="progress glass a d1">
      <div class="progress-info">
        <span class="gold-text">{{ gacha.progress?.pool ?? '准备中' }}</span>
        <span class="p-text">{{ gacha.progress ? `${gacha.progress.index} / ${gacha.progress.total}` : '连接日志…' }}</span>
      </div>
      <div class="track">
        <div class="fill" :style="{ width: gacha.progress ? `${(gacha.progress.index / gacha.progress.total) * 100}%` : '6%' }" />
      </div>
    </div>
    <div v-else-if="gacha.lastError" class="error glass a d1">{{ gacha.lastError }}</div>

    <!-- 空状态 -->
    <EmptyState
      v-if="!activeDetail"
      icon="M3 3v18h18M7 15l4-6 4 3 5-8"
      title="暂无唤取数据"
      description="启动《鸣潮》并打开「唤取 → 唤取记录」页面，然后点击「立即抓取」。aMC 会自动发现客户端日志、提取凭证并从官方接口拉取全部 13 种卡池记录。"
    />

    <!-- 数据面板 -->
    <template v-else>
      <div class="tabs glass a d2">
        <div
          v-for="tab in poolTabs"
          :key="tab.key"
          class="tab"
          :class="{ active: activePool === tab.key }"
          @click="activePool = tab.key"
        >
          {{ tab.name }}<span class="n">{{ tab.count }}</span>
        </div>
      </div>

      <section class="stats">
        <div class="stat glass a d3">
          <div class="halo" />
          <div class="k">累计唤取</div>
          <div class="v gold-text">{{ activeStats?.total.toLocaleString() }}</div>
          <div class="foot">当前筛选 <b>{{ activePoolName }}</b></div>
        </div>
        <div class="stat glass a d3">
          <div class="k">5★ 已获取</div>
          <div class="v">{{ activeStats?.count5 }}<em>★</em></div>
          <div class="foot">4★ 共 {{ activeStats?.count4 }} 个</div>
        </div>
        <div class="stat glass a d4">
          <div class="k">平均出金</div>
          <div class="v">{{ activeStats?.avgPity ? activeStats.avgPity.toFixed(1) : '—' }}<em>抽</em></div>
          <div class="foot">综合期望 62.0 抽</div>
        </div>
        <div class="stat glass a d4">
          <div class="k">欧非指数</div>
          <div class="v">{{ activeStats?.luckIndex ? activeStats.luckIndex.toFixed(2) : '—' }}</div>
          <div class="foot">较均值<b>{{ luckLabel }}</b></div>
        </div>
      </section>

      <section class="mid">
        <div class="card glass a d4">
          <h3>保底进度<span class="tag">{{ activePoolName }}</span></h3>
          <template v-if="gacha.activePoolStats">
            <div class="pity">
              <b class="gold-text">{{ gacha.activePoolStats.pity }}</b><span>/ 80</span>
              <span class="left">距保底还有 <b>{{ Math.max(80 - gacha.activePoolStats.pity, 0) }}</b> 抽</span>
            </div>
            <div class="track"><div class="fill" :style="{ width: `${gacha.pityRatio * 100}%` }" /></div>
            <div class="pity-fn">
              <span>{{ gacha.activePoolStats.pityIsFloor ? '窗口内无 5★，数值为下界' : '软保底 65 抽起概率提升' }}</span>
              <span>硬保底 80 抽</span>
            </div>
          </template>
          <template v-else>
            <div class="pity"><span class="allhint">选择单个卡池查看保底进度</span></div>
            <div class="track"><div class="fill" style="width: 0%" /></div>
            <div class="pity-fn"><span>—</span><span>硬保底 80 抽</span></div>
          </template>
        </div>

        <div class="card glass a d5">
          <h3>账号概览<span class="tag">ACCOUNT</span></h3>
          <div class="overview">
            <div class="ov-row"><span>UID</span><b>{{ activeDetail.uid }}</b></div>
            <div class="ov-row"><span>区服</span><b>{{ activeDetail.svrArea === 'global' ? '国际服' : '国服' }}</b></div>
            <div class="ov-row"><span>最近更新</span><b>{{ activeDetail.fetchedAt || '—' }}</b></div>
            <div class="ov-row"><span>5★ 构成</span><b>限定 {{ star5Split.limited }} · 常驻 {{ star5Split.standard }}</b></div>
          </div>
          <div class="bnr-foot">
            <span class="chip">{{ activeStats?.pools.length ?? 0 }} 个卡池有记录</span>
            <span class="chip dim">5★ ×{{ activeStats?.count5 }}</span>
          </div>
        </div>

        <div class="card glass a d5">
          <h3>5★ 构成<span class="tag">RATIO</span></h3>
          <div class="donut-wrap">
            <div
              class="donut"
              :style="{ background: `conic-gradient(#EDBE5A 0 ${limitedPct}%, rgba(255,255,255,.16) ${limitedPct}% 100%)` }"
            />
            <div class="legend">
              <div class="l1"><i />限定池<b>{{ star5Split.limited }}</b><small>{{ limitedPct }}%</small></div>
              <div class="l2"><i />常驻类<b>{{ star5Split.standard }}</b><small>{{ 100 - limitedPct }}%</small></div>
            </div>
          </div>
        </div>
      </section>

      <section class="bottom">
        <div class="card glass a d5">
          <h3>最近 5★ 记录<span class="tag">RECENT</span></h3>
          <table>
            <thead>
              <tr>
                <th style="width: 84px">日期</th>
                <th>名称</th>
                <th style="width: 140px">卡池</th>
                <th style="width: 90px">间隔</th>
                <th style="width: 84px">结果</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in gacha.recent5Filtered" :key="r.time + r.name">
                <td class="mono">{{ fmtDate(r.time) }}</td>
                <td><span class="who">{{ r.name }}</span></td>
                <td class="pool">{{ poolNameOf(r.pool) }}</td>
                <td class="mono">{{ r.gap >= 0 ? `${r.gap} 抽` : '—' }}</td>
                <td>
                  <span v-if="r.gap > 0 && r.gap <= 20" class="pill p-off">欧皇</span>
                  <span v-else-if="r.gap >= 0" class="pill p-gold">不歪</span>
                  <span v-else class="pill p-lost">窗口首金</span>
                </td>
              </tr>
              <tr v-if="!gacha.recent5Filtered.length">
                <td colspan="5" class="none">当前筛选下暂无 5★ 记录</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="card glass a d6">
          <h3>出金分布<span class="tag">DISTRIBUTION</span></h3>
          <div class="bars">
            <div
              v-for="(v, i) in gacha.distFiltered"
              :key="i"
              class="bar"
              :class="{ hot: v > 0 && v === distMax }"
            >
              <em>{{ v }}</em>
              <i :style="{ height: `${(v / distMax) * 100}%`, minHeight: v > 0 ? '3px' : '2px' }" />
              <span>{{ distLabel(i) }}</span>
            </div>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 14px;
  height: 100%;
  overflow: hidden;
}
/* 账号切换 */
.account {
  position: relative;
  display: flex;
  align-items: center;
  gap: 9px;
  height: 38px;
  padding: 0 15px;
  border-radius: 999px;
  font-size: 12.5px;
  color: var(--txt-2);
  cursor: pointer;
  user-select: none;
}
.account b {
  color: var(--txt-1);
  font-weight: 600;
}
.caret {
  width: 8px;
  height: 8px;
  border-right: 1.5px solid var(--txt-3);
  border-bottom: 1.5px solid var(--txt-3);
  transform: rotate(45deg) translateY(-2px);
}
.menu {
  position: absolute;
  top: 46px;
  right: 0;
  min-width: 210px;
  padding: 6px;
  z-index: 30;
}
.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 14px;
  padding: 9px 12px;
  border-radius: var(--radius-xs);
  font-size: 12.5px;
  cursor: pointer;
}
.menu-item:hover {
  background: rgba(255, 255, 255, 0.06);
}
.menu-item.on {
  color: var(--warn);
}
.menu-item small {
  color: var(--txt-3);
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
  font-family: var(--font-ui);
}
.btn:hover {
  transform: translateY(-1px);
}
.btn:disabled {
  opacity: 0.7;
  cursor: wait;
  transform: none;
}
/* 抓取进度 */
.progress {
  padding: 13px 20px;
  flex-shrink: 0;
}
.progress-info {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  margin-bottom: 9px;
}
.p-text {
  color: var(--txt-2);
}
.track {
  height: 7px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.07);
  overflow: hidden;
}
.fill {
  height: 100%;
  border-radius: 5px;
  background: linear-gradient(90deg, var(--gold-3), var(--gold-1));
  box-shadow: 0 0 14px rgba(244, 206, 106, 0.5);
  transition: width 0.45s cubic-bezier(0.2, 0.7, 0.2, 1);
}
.error {
  padding: 13px 20px;
  font-size: 13px;
  color: var(--danger);
  flex-shrink: 0;
}
/* tabs */
.tabs {
  display: flex;
  gap: 6px;
  padding: 8px 10px;
  width: fit-content;
  flex-shrink: 0;
  max-width: 100%;
  overflow-x: auto;
}
.tab {
  height: 34px;
  padding: 0 17px;
  display: flex;
  align-items: center;
  gap: 7px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--txt-2);
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
  white-space: nowrap;
  flex-shrink: 0;
}
.tab:hover {
  color: var(--txt-1);
}
.tab.active {
  background: rgba(255, 255, 255, 0.08);
  color: var(--txt-1);
  border-color: var(--line);
  font-weight: 600;
}
.tab .n {
  font-size: 11px;
  color: var(--txt-3);
}
.tab.active .n {
  color: var(--warn);
}
/* stats */
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
  letter-spacing: 0;
}
.stat .foot {
  margin-top: 9px;
  font-size: 12px;
  color: var(--txt-2);
}
.stat .foot b {
  color: var(--ok);
  font-weight: 600;
}
/* mid */
.mid {
  display: grid;
  grid-template-columns: 1.25fr 1.15fr 1fr;
  gap: 14px;
  flex-shrink: 0;
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
}
.card h3 .tag {
  font-size: 10.5px;
  color: var(--txt-3);
  font-weight: 500;
  letter-spacing: 0.08em;
}
.pity {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin: 16px 0 12px;
}
.pity b {
  font-size: 44px;
  font-weight: 650;
  letter-spacing: -1.5px;
}
.pity span {
  font-size: 16px;
  color: var(--txt-3);
  font-weight: 500;
}
.pity .left {
  margin-left: auto;
  font-size: 12.5px;
  color: var(--txt-2);
}
.pity .left b {
  color: var(--warn);
  font-size: 15px;
}
.allhint {
  font-size: 13.5px;
  color: var(--txt-2);
  margin: 14px 0;
}
.pity-fn {
  display: flex;
  justify-content: space-between;
  font-size: 11.5px;
  color: var(--txt-3);
  margin-top: 10px;
}
/* account overview */
.overview {
  margin-top: 14px;
  display: flex;
  flex-direction: column;
  gap: 11px;
  flex: 1;
}
.ov-row {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  color: var(--txt-3);
}
.ov-row b {
  color: var(--txt-1);
  font-weight: 600;
}
.bnr-foot {
  display: flex;
  gap: 8px;
  padding-top: 14px;
}
.chip {
  height: 26px;
  padding: 0 12px;
  border-radius: 999px;
  font-size: 11.5px;
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--gold-soft);
  color: var(--warn);
  border: 1px solid rgba(237, 190, 90, 0.2);
}
.chip.dim {
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-2);
  border-color: var(--line);
}
/* donut */
.donut-wrap {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
  margin-top: 12px;
}
.donut {
  width: 108px;
  height: 108px;
  border-radius: 50%;
  flex-shrink: 0;
  -webkit-mask: radial-gradient(circle, transparent 56%, #000 57.5%);
  mask: radial-gradient(circle, transparent 56%, #000 57.5%);
}
.legend {
  display: flex;
  flex-direction: column;
  gap: 11px;
  font-size: 12.5px;
}
.legend i {
  width: 9px;
  height: 9px;
  border-radius: 3px;
  display: inline-block;
  margin-right: 8px;
}
.legend b {
  font-size: 16px;
  font-weight: 650;
  margin-left: 6px;
}
.legend .l1 i {
  background: #edbe5a;
}
.legend .l2 i {
  background: rgba(255, 255, 255, 0.35);
}
.legend small {
  color: var(--txt-3);
  margin-left: 6px;
  font-size: 11px;
}
/* bottom */
.bottom {
  display: grid;
  grid-template-columns: 1.55fr 1fr;
  gap: 14px;
  flex: 1;
  min-height: 0;
}
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 8px;
}
th {
  font-size: 11px;
  color: var(--txt-3);
  font-weight: 500;
  text-align: left;
  padding: 10px 8px 8px;
  letter-spacing: 0.1em;
  border-bottom: 1px solid var(--line-2);
}
td {
  padding: 10px 8px;
  font-size: 13.5px;
  border-bottom: 1px solid var(--line-2);
}
tr:last-child td {
  border-bottom: none;
}
td .who {
  font-weight: 600;
}
td .pool {
  color: var(--txt-3);
  font-size: 12px;
}
td.mono {
  font-family: var(--font-mono);
  font-size: 12.5px;
  color: var(--txt-2);
}
td.none {
  text-align: center;
  color: var(--txt-3);
  font-size: 12.5px;
  padding: 22px 8px;
}
.pill {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 11px;
  border-radius: 999px;
  font-size: 11.5px;
  font-weight: 600;
}
.p-gold {
  background: var(--gold-soft);
  color: var(--warn);
  border: 1px solid rgba(237, 190, 90, 0.2);
}
.p-off {
  background: rgba(125, 224, 168, 0.12);
  color: var(--ok);
  border: 1px solid rgba(125, 224, 168, 0.2);
}
.p-lost {
  background: rgba(167, 139, 250, 0.13);
  color: var(--info);
  border: 1px solid rgba(167, 139, 250, 0.22);
}
/* bars */
.bars {
  display: flex;
  align-items: flex-end;
  gap: 9px;
  height: 150px;
  margin-top: auto;
  padding-top: 12px;
}
.bar {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 7px;
  height: 100%;
  justify-content: flex-end;
}
.bar i {
  width: 100%;
  max-width: 26px;
  border-radius: 6px 6px 3px 3px;
  background: rgba(255, 255, 255, 0.12);
  transition: height 0.5s cubic-bezier(0.2, 0.7, 0.2, 1);
}
.bar.hot i {
  background: linear-gradient(180deg, var(--gold-1), var(--gold-3));
  box-shadow: 0 0 16px rgba(244, 206, 106, 0.35);
}
.bar span {
  font-size: 10px;
  color: var(--txt-3);
  white-space: nowrap;
}
.bar em {
  font-style: normal;
  font-size: 11px;
  color: var(--txt-2);
  font-weight: 600;
}
</style>
