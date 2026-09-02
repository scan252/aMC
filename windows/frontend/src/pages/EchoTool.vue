<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { GameDataService } = bindings as any

interface SubstatDef { name: string; max: number }
interface Scheme { name: string; note: string }

const substats = ref<SubstatDef[]>([])
const schemes = ref<Scheme[]>([])
const scheme = ref('')
const loading = ref(true)

const slots = ref(
  Array.from({ length: 5 }, () => ({ name: '', value: null as number | null })),
)

const result = ref<{
  equivalentRolls: number
  percent: number
  grade: string
  detail: { name: string; value: number; rolls: number; weight: number }[]
  unknown?: string[]
} | null>(null)

async function loadMeta() {
  const [tables, _echoes] = ((await GameDataService.EchoScoreMeta()) as any[]) ?? [{}, []]
  substats.value = tables?.substats ?? []
  schemes.value = tables?.schemes ?? []
  if (schemes.value.length) {
    scheme.value = schemes.value[0].name
  }
  loading.value = false
}

async function score() {
  const subs = slots.value
    .filter((s) => s.name && s.value !== null && !Number.isNaN(s.value))
    .map((s) => ({ name: s.name, value: Number(s.value) }))
  if (!subs.length) {
    result.value = null
    return
  }
  result.value = ((await GameDataService.ScoreEcho(scheme.value, subs)) as any) ?? null
}

const gradeColor = computed(() => {
  const g = result.value?.grade
  return { S: '#F4CE6A', A: '#7DE0A8', B: '#A8C8F7', C: '#BCA8F7', D: '#8A8A94' }[g ?? 'D'] ?? '#8A8A94'
})

function namesFor(i: number) {
  // 同一副词条不允许重复选择
  return substats.value.filter((s) => !slots.value.some((o, j) => j !== i && o.name === s.name))
}

onMounted(loadMeta)
</script>

<template>
  <div class="page">
    <PageHeader title="声骸工具" subtitle="词条评分 · 等效满-roll 词条 · 权重方案可切换" />

    <div class="layout">
      <div class="editor glass a d2">
        <h3>副词条录入<span class="tag">INPUT</span></h3>
        <div class="slot" v-for="(slot, i) in slots" :key="i">
          <span class="idx">{{ i + 1 }}</span>
          <select v-model="slot.name" class="field" @change="score">
            <option value="">选择词条</option>
            <option v-for="s in namesFor(i)" :key="s.name" :value="s.name">{{ s.name }}</option>
          </select>
          <input
            v-model.number="slot.value"
            class="field val"
            type="number"
            step="0.1"
            min="0"
            placeholder="数值"
            :disabled="!slot.name"
            @input="score"
          />
          <span class="max" v-if="slot.name">/ {{ substats.find((s) => s.name === slot.name)?.max }}</span>
        </div>

        <div class="scheme-row">
          <span class="lbl">权重方案</span>
          <select v-model="scheme" class="field" @change="score">
            <option v-for="sc in schemes" :key="sc.name" :value="sc.name">{{ sc.name }}</option>
          </select>
        </div>
        <p class="note">{{ schemes.find((s) => s.name === scheme)?.note }}</p>
      </div>

      <div class="result glass a d3">
        <h3>评分结果<span class="tag">SCORE</span></h3>
        <template v-if="result">
          <div class="hero">
            <div class="grade" :style="{ color: gradeColor, borderColor: `${gradeColor}66` }">{{ result.grade }}</div>
            <div class="nums">
              <div class="num-big"><b class="gold-text">{{ result.percent.toFixed(1) }}</b><span>%</span></div>
              <div class="rolls">等效满-roll 词条 <b>{{ result.equivalentRolls.toFixed(2) }}</b> / 10</div>
            </div>
          </div>
          <div class="detail">
            <div v-for="d in result.detail" :key="d.name" class="d-row">
              <span class="d-name" :class="{ off: d.weight === 0 }">{{ d.name }}</span>
              <div class="d-bar"><i :style="{ width: `${Math.min(d.rolls * 20, 100)}%` }" /></div>
              <span class="d-rolls">{{ d.rolls.toFixed(1) }}</span>
            </div>
          </div>
          <p v-if="result.unknown?.length" class="unknown">已忽略未知词条：{{ result.unknown.join('、') }}</p>
        </template>
        <div v-else class="empty-hint">在左侧录入副词条后自动评分。<br />词条满值与权重方案可在数据版本通道中校准。</div>
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
  overflow-y: auto;
}
.layout {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 14px;
  align-items: start;
}
.editor,
.result {
  padding: 20px 22px;
}
h3 {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--txt-2);
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}
h3 .tag {
  font-size: 10.5px;
  color: var(--txt-3);
  letter-spacing: 0.08em;
}
.slot {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 9px;
}
.idx {
  width: 22px;
  height: 22px;
  border-radius: 7px;
  display: grid;
  place-items: center;
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-3);
  font-size: 11px;
  flex-shrink: 0;
}
.field {
  height: 34px;
  border-radius: var(--radius-xs);
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-1);
  padding: 0 11px;
  font-size: 12.5px;
  outline: none;
  font-family: var(--font-ui);
  flex: 1;
  min-width: 0;
}
.field:focus {
  border-color: rgba(237, 190, 90, 0.4);
}
.field.val {
  max-width: 90px;
  flex: 0 0 90px;
}
.max {
  font-size: 10.5px;
  color: var(--txt-3);
  font-family: var(--font-mono);
  flex-shrink: 0;
}
.scheme-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 16px;
}
.lbl {
  font-size: 12.5px;
  color: var(--txt-2);
}
.note {
  margin-top: 9px;
  font-size: 11.5px;
  color: var(--txt-3);
  line-height: 1.6;
}
.hero {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 18px;
}
.grade {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: grid;
  place-items: center;
  font-size: 34px;
  font-weight: 750;
  border: 1px solid;
  background: rgba(255, 255, 255, 0.03);
}
.num-big b {
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -1px;
}
.num-big span {
  color: var(--txt-3);
  font-size: 16px;
  margin-left: 3px;
}
.rolls {
  font-size: 12px;
  color: var(--txt-2);
  margin-top: 4px;
}
.rolls b {
  color: var(--txt-1);
}
.detail {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.d-row {
  display: grid;
  grid-template-columns: 110px 1fr 42px;
  align-items: center;
  gap: 10px;
}
.d-name {
  font-size: 12px;
  color: var(--txt-1);
}
.d-name.off {
  color: var(--txt-3);
  text-decoration: line-through;
}
.d-bar {
  height: 7px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.07);
  overflow: hidden;
}
.d-bar i {
  display: block;
  height: 100%;
  border-radius: 4px;
  background: linear-gradient(90deg, var(--gold-3), var(--gold-1));
  box-shadow: 0 0 10px rgba(244, 206, 106, 0.4);
  transition: width 0.4s ease;
}
.d-rolls {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--txt-2);
  text-align: right;
}
.unknown {
  margin-top: 12px;
  font-size: 11.5px;
  color: var(--txt-3);
}
.empty-hint {
  text-align: center;
  color: var(--txt-3);
  font-size: 12.5px;
  line-height: 2;
  padding: 40px 10px;
}
</style>
