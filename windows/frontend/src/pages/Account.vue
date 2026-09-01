<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import { useAccountStore } from '../stores/account'

const account = useAccountStore()

const phone = ref('')
const code = ref('')
const step = ref<'phone' | 'code'>('phone')

async function sendSms() {
  await account.sendSms(phone.value.trim())
  if (!account.lastError) step.value = 'code'
}

async function login() {
  await account.login(phone.value.trim(), code.value.trim())
  step.value = 'phone'
  phone.value = ''
  code.value = ''
}

onMounted(() => account.load())
</script>

<template>
  <div class="page">
    <PageHeader title="账号中心" subtitle="库街区账号绑定 · 每日签到 · 结晶波片">
      <template #actions>
        <span v-if="account.overview?.status.mode === 'mock'" class="mode-chip">演示模式</span>
        <button v-if="account.bound" class="btn ghost" @click="account.logout()">解绑</button>
      </template>
    </PageHeader>

    <div v-if="account.lastError" class="msg error glass">{{ account.lastError }}</div>
    <div v-else-if="account.lastMessage" class="msg ok glass">{{ account.lastMessage }}</div>

    <!-- 未绑定：登录流程 -->
    <div v-if="!account.bound" class="login glass a d2">
      <div class="login-inner">
        <div class="icon-wrap">
          <svg viewBox="0 0 24 24">
            <path d="M12 3a4 4 0 1 1 0 8 4 4 0 0 1 0-8ZM4 21c1.5-4 4.5-6 8-6s6.5 2 8 6" />
          </svg>
        </div>
        <h2>绑定库街区账号</h2>
        <p class="hint">使用手机号 + 短信验证码登录，全程走库街区官方接口。当前为演示模式，验证码固定为 <b>888888</b>。</p>

        <div class="form">
          <input v-model="phone" class="field" type="text" maxlength="11" placeholder="手机号" :disabled="step === 'code'" />
          <template v-if="step === 'code'">
            <input v-model="code" class="field" type="text" maxlength="6" placeholder="验证码" @keyup.enter="login" />
            <button class="btn" :disabled="account.loading" @click="login">登录</button>
            <button class="btn ghost" @click="step = 'phone'">返回</button>
          </template>
          <template v-else>
            <button class="btn" :disabled="account.loading || phone.length !== 11" @click="sendSms">获取验证码</button>
          </template>
        </div>
      </div>
    </div>

    <!-- 已绑定：概览 -->
    <template v-else>
      <section class="cards">
        <div class="card glass a d2">
          <h3>每日签到<span class="tag">SIGN-IN</span></h3>
          <div class="sign-body">
            <div class="sign-state" :class="{ done: account.overview?.signIn?.hadSignIn }">
              {{ account.overview?.signIn?.hadSignIn ? '今日已签到' : '今日未签到' }}
            </div>
            <div class="sign-meta">本月累计 <b>{{ account.overview?.signIn?.totalSignIn ?? 0 }}</b> 天 · 今日奖励 {{ account.overview?.signIn?.todayReward || '—' }}</div>
            <button class="btn" :disabled="account.overview?.signIn?.hadSignIn || account.loading" @click="account.signInNow()">
              {{ account.overview?.signIn?.hadSignIn ? '已完成' : '立即签到' }}
            </button>
          </div>
        </div>

        <div class="card glass a d3">
          <h3>结晶波片<span class="tag">WAVEPLATE</span></h3>
          <div class="wave-body">
            <div class="wave-num">
              <b class="gold-text">{{ account.waveplate?.cur ?? '—' }}</b>
              <span>/ {{ account.waveplate?.max ?? 240 }}</span>
            </div>
            <div class="track">
              <div class="fill" :style="{ width: `${((account.waveplate?.cur ?? 0) / (account.waveplate?.max ?? 240)) * 100}%` }" />
            </div>
            <div class="wave-meta">回满约需 <b>{{ account.waveFullInMin }}</b> 分钟 · 领取后自动提醒</div>
          </div>
        </div>

        <div class="card glass a d4">
          <h3>绑定角色<span class="tag">ROLES</span></h3>
          <div class="roles">
            <div v-for="role in account.overview?.roles ?? []" :key="role.roleId" class="role-row">
              <div>
                <b>{{ role.roleName }}</b>
                <small>{{ role.areaName }} · Lv.{{ role.level }}</small>
              </div>
              <span v-if="role.isDefault" class="chip">默认</span>
            </div>
            <div v-if="!account.overview?.roles?.length" class="roles-empty">暂无绑定角色</div>
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
  overflow-y: auto;
}
.mode-chip {
  height: 26px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border-radius: 999px;
  font-size: 11.5px;
  background: var(--gold-soft);
  color: var(--warn);
  border: 1px solid rgba(237, 190, 90, 0.2);
}
.msg {
  padding: 12px 20px;
  font-size: 13px;
  border-radius: var(--radius);
  flex-shrink: 0;
}
.msg.error {
  color: var(--danger);
}
.msg.ok {
  color: var(--ok);
}
.login {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 0;
}
.login-inner {
  text-align: center;
  max-width: 420px;
  padding: 30px;
}
.icon-wrap {
  width: 74px;
  height: 74px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: var(--gold-soft);
  border: 1px solid rgba(237, 190, 90, 0.2);
  margin: 0 auto 20px;
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
.hint {
  font-size: 12.5px;
  color: var(--txt-3);
  line-height: 1.8;
  margin-bottom: 22px;
}
.hint b {
  color: var(--warn);
}
.form {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}
.field {
  height: 38px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.05);
  color: var(--txt-1);
  padding: 0 18px;
  font-size: 13.5px;
  outline: none;
  width: 200px;
  font-family: var(--font-ui);
}
.field:focus {
  border-color: rgba(237, 190, 90, 0.4);
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
  box-shadow: var(--gold-glow), inset 0 1px 0 rgba(255, 255, 255, 0.5);
  font-family: var(--font-ui);
  transition: all 0.2s;
}
.btn:hover {
  transform: translateY(-1px);
}
.btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
  transform: none;
}
.btn.ghost {
  background: rgba(255, 255, 255, 0.06);
  color: var(--txt-2);
  box-shadow: none;
  border: 1px solid var(--line);
}
.cards {
  display: grid;
  grid-template-columns: 1.2fr 1.2fr 1fr;
  gap: 14px;
}
.card {
  padding: 20px 22px;
  display: flex;
  flex-direction: column;
  min-height: 190px;
}
.card h3 {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--txt-2);
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}
.card h3 .tag {
  font-size: 10.5px;
  color: var(--txt-3);
  letter-spacing: 0.08em;
}
.sign-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}
.sign-state {
  font-size: 22px;
  font-weight: 650;
}
.sign-state.done {
  color: var(--ok);
}
.sign-meta {
  font-size: 12.5px;
  color: var(--txt-3);
}
.sign-meta b {
  color: var(--txt-1);
}
.sign-body .btn {
  margin-top: auto;
  align-self: flex-start;
}
.wave-num {
  display: flex;
  align-items: baseline;
  gap: 5px;
  margin-bottom: 12px;
}
.wave-num b {
  font-size: 40px;
  font-weight: 650;
  letter-spacing: -1px;
}
.wave-num span {
  color: var(--txt-3);
  font-size: 15px;
}
.track {
  height: 8px;
  border-radius: 5px;
  background: rgba(255, 255, 255, 0.07);
  overflow: hidden;
  margin-bottom: 10px;
}
.fill {
  height: 100%;
  border-radius: 5px;
  background: linear-gradient(90deg, var(--gold-3), var(--gold-1));
  box-shadow: 0 0 14px rgba(244, 206, 106, 0.5);
  transition: width 0.6s cubic-bezier(0.2, 0.7, 0.2, 1);
}
.wave-meta {
  font-size: 12.5px;
  color: var(--txt-3);
}
.wave-meta b {
  color: var(--warn);
}
.roles {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}
.role-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 13px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--line-2);
}
.role-row b {
  display: block;
  font-size: 13.5px;
}
.role-row small {
  color: var(--txt-3);
  font-size: 11.5px;
}
.role-row .chip {
  height: 24px;
  padding: 0 11px;
  border-radius: 999px;
  font-size: 11px;
  display: flex;
  align-items: center;
  background: var(--gold-soft);
  color: var(--warn);
}
.roles-empty {
  color: var(--txt-3);
  font-size: 12.5px;
  text-align: center;
  padding: 20px;
}
</style>
