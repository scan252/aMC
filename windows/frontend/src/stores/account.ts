import { defineStore } from 'pinia'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { KurobbsService } = bindings as any

export interface KurobbsStatus {
  bound: boolean
  name: string
  phone: string
  mode: string
  userId: string
}

export interface Role {
  roleId: string
  serverId: string
  roleName: string
  level: number
  areaName: string
  isDefault: boolean
}

export interface SignInInfo {
  hadSignIn: boolean
  monthStart: string
  totalSignIn: number
  todayReward: string
}

export interface WidgetData {
  energy: { cur: number; max: number; Full?: string }
  chest: number
  combat: string
  unlock: string
}

export interface KurobbsOverview {
  status: KurobbsStatus
  roles: Role[]
  signIn: SignInInfo | null
  widget: WidgetData | null
}

export const useAccountStore = defineStore('account', {
  state: () => ({
    overview: null as KurobbsOverview | null,
    loading: false as boolean,
    lastError: '' as string,
    lastMessage: '' as string,
  }),

  getters: {
    bound: (state) => state.overview?.status.bound ?? false,
    waveplate(state): { cur: number; max: number } | null {
      const w = state.overview?.widget
      if (!w) return null
      return { cur: w.energy?.cur ?? 0, max: w.energy?.max ?? 240 }
    },
    waveFullInMin(state): number {
      const w = state.overview?.widget
      if (!w?.energy) return 0
      const full = (w.energy as any).Full ?? (w.energy as any).full
      if (typeof full === 'string' && full) {
        const ms = new Date(full).getTime() - Date.now()
        return Math.max(0, Math.round(ms / 60000))
      }
      return Math.max(0, (w.energy.max - w.energy.cur) * 6)
    },
  },

  actions: {
    async load() {
      this.loading = true
      this.lastError = ''
      try {
        this.overview = ((await KurobbsService.GetOverview()) as any) ?? null
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      } finally {
        this.loading = false
      }
    },

    async sendSms(phone: string) {
      this.lastError = ''
      this.lastMessage = ''
      try {
        await KurobbsService.SendSms(phone)
        this.lastMessage = '验证码已发送（演示模式：验证码为 888888）'
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      }
    },

    async login(phone: string, code: string) {
      this.loading = true
      this.lastError = ''
      try {
        this.overview = ((await KurobbsService.Login(phone, code)) as any) ?? null
        await this.load()
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      } finally {
        this.loading = false
      }
    },

    async logout() {
      await KurobbsService.Logout()
      await this.load()
    },

    async signInNow() {
      this.lastError = ''
      this.lastMessage = ''
      try {
        const msg = (await KurobbsService.SignInNow()) as unknown as string
        this.lastMessage = msg
        await this.load()
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      }
    },
  },
})
