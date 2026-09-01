<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import { useSettingsStore } from '../stores/settings'

const settings = useSettingsStore()
const draft = ref({
  autostart: false,
  signInAuto: true,
  signInHour: 8,
  waveNotify: true,
  logPath: '',
  language: 'zh-CN',
})

watch(
  () => settings.settings,
  (v) => {
    if (v) {
      draft.value = {
        autostart: v.autostart,
        signInAuto: v.signInAuto,
        signInHour: v.signInHour,
        waveNotify: v.waveNotify,
        logPath: v.logPath,
        language: v.language,
      }
    }
  },
  { immediate: true },
)

onMounted(() => settings.load())
</script>

<template>
  <div class="page">
    <PageHeader title="设置" :subtitle="settings.savedFlash ? '✓ 已保存' : '应用行为与数据管理'">
      <template #actions>
        <button class="btn" :disabled="settings.saving" @click="settings.save()">
          {{ settings.saving ? '保存中…' : '保存设置' }}
        </button>
      </template>
    </PageHeader>

    <div v-if="settings.lastError" class="msg error glass">{{ settings.lastError }}</div>

    <div v-if="settings.settings" class="sections">
      <section class="sec glass a d2">
        <h3>启动与常驻<span class="tag">GENERAL</span></h3>
        <label class="row">
          <div class="row-text">
            <b>开机自启</b>
            <small>登录 Windows 后自动运行 aMC Suite（当前注册表状态：{{ settings.settings.autostartOn ? '已启用' : '未启用' }}）</small>
          </div>
          <input v-model="draft.autostart" type="checkbox" class="switch" />
        </label>
      </section>

      <section class="sec glass a d3">
        <h3>账号任务<span class="tag">TASKS</span></h3>
        <label class="row">
          <div class="row-text">
            <b>每日自动签到</b>
            <small>绑定库街区账号后，在指定时刻自动执行签到</small>
          </div>
          <input v-model="draft.signInAuto" type="checkbox" class="switch" />
        </label>
        <label class="row">
          <div class="row-text">
            <b>自动签到时刻</b>
            <small>每天该时刻后第一次检查时执行</small>
          </div>
          <select v-model.number="draft.signInHour" class="select">
            <option v-for="h in 24" :key="h - 1" :value="h - 1">{{ String(h - 1).padStart(2, '0') }}:00</option>
          </select>
        </label>
        <label class="row">
          <div class="row-text">
            <b>波片回满提醒</b>
            <small>结晶波片回满时弹出通知（约每 6 小时最多一次）</small>
          </div>
          <input v-model="draft.waveNotify" type="checkbox" class="switch" />
        </label>
      </section>

      <section class="sec glass a d4">
        <h3>数据<span class="tag">DATA</span></h3>
        <div class="row static">
          <div class="row-text">
            <b>数据目录</b>
            <small>抽卡记录与账号数据保存位置（与 Mac 版 aMC 格式互通）</small>
          </div>
          <code class="path">{{ settings.settings.dataDir }}</code>
        </div>
        <label class="row">
          <div class="row-text">
            <b>手动指定游戏日志</b>
            <small>自动发现失败时，可填写 Client.log 完整路径（留空使用自动发现）</small>
          </div>
          <input v-model="draft.logPath" class="input" type="text" placeholder="D:\...\Wuthering Waves Game\Client\Saved\Logs\Client.log" />
        </label>
      </section>
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
.btn:disabled {
  opacity: 0.6;
}
.msg {
  padding: 12px 20px;
  font-size: 13px;
  border-radius: var(--radius);
  color: var(--danger);
  flex-shrink: 0;
}
.sections {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.sec {
  padding: 6px 22px;
}
.sec h3 {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--txt-2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0 4px;
  border-bottom: 1px solid var(--line-2);
}
.sec h3 .tag {
  font-size: 10.5px;
  color: var(--txt-3);
  letter-spacing: 0.08em;
}
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 15px 0;
  border-bottom: 1px solid var(--line-2);
  cursor: pointer;
}
.row:last-child {
  border-bottom: none;
}
.row.static {
  cursor: default;
}
.row-text b {
  display: block;
  font-size: 13.5px;
  font-weight: 600;
}
.row-text small {
  display: block;
  margin-top: 4px;
  font-size: 11.5px;
  color: var(--txt-3);
  line-height: 1.6;
}
.switch {
  appearance: none;
  width: 42px;
  height: 24px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
  position: relative;
  cursor: pointer;
  transition: background 0.2s;
  flex-shrink: 0;
}
.switch::after {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.7);
  transition: all 0.2s;
}
.switch:checked {
  background: linear-gradient(135deg, var(--gold-1), var(--gold-2));
}
.switch:checked::after {
  left: 21px;
  background: #201503;
}
.select,
.input {
  height: 34px;
  border-radius: var(--radius-xs);
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-1);
  padding: 0 12px;
  font-size: 12.5px;
  outline: none;
  font-family: var(--font-ui);
  flex-shrink: 0;
}
.input {
  width: 340px;
}
.input:focus {
  border-color: rgba(237, 190, 90, 0.4);
}
.path {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--txt-2);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--line-2);
  border-radius: var(--radius-xs);
  padding: 8px 12px;
  max-width: 380px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
