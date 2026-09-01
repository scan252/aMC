import { defineStore } from 'pinia'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'
import { Events } from '@wailsio/runtime'

const { GachaService } = bindings

export interface AccountSummary {
  uid: string
  svrArea: string
  fetchedAt: string
  total: number
  count5: number
  avgPity: number
}

export interface GachaRecord {
  cardPoolType: string
  resourceId: number
  qualityLevel: number
  resourceType: string
  name: string
  count: number
  time: string
}

export interface PoolStats {
  poolType: number
  poolName: string
  total: number
  count5: number
  count4: number
  count3: number
  avgPity: number
  pity: number
  pityIsFloor: boolean
  last5Name: string
  last5Time: string
}

export interface OverallStats {
  total: number
  count5: number
  count4: number
  avgPity: number
  luckIndex: number
  pools: PoolStats[]
}

export interface RecentItem extends GachaRecord {
  pool: string
  gap: number
}

export interface AccountDetail {
  uid: string
  svrArea: string
  fetchedAt: string
  stats: OverallStats
  dist: Record<string, number[]>
  recent5: RecentItem[]
}

export interface FetchProgress {
  index: number
  total: number
  pool: string
  err?: string
}

export const POOL_NAMES: Record<string, string> = {
  '1': '角色活动唤取', '2': '武器活动唤取', '3': '角色常驻唤取', '4': '武器常驻唤取',
  '5': '新手唤取', '6': '新手自选唤取', '7': '感恩定向唤取', '8': '角色新旅唤取',
  '9': '武器新旅唤取', '10': '角色联动唤取', '11': '武器联动唤取', '12': '角色忆旅唤取',
  '13': '武器忆旅唤取',
}

export const useGachaStore = defineStore('gacha', {
  state: () => ({
    accounts: [] as AccountSummary[],
    detail: null as AccountDetail | null,
    activePool: 'all' as string,
    loading: false as boolean,
    progress: null as FetchProgress | null,
    lastError: '' as string,
  }),

  getters: {
    /** 当前激活卡池的统计 */
    activePoolStats(state): PoolStats | null {
      if (!state.detail) return null
      if (state.activePool === 'all') return null
      return state.detail.stats.pools.find((p) => String(p.poolType) === state.activePool) ?? null
    },
    /** 保底进度：0..1 */
    pityRatio(state): number {
      if (!state.detail || state.activePool === 'all') return 0
      const p = state.detail.stats.pools.find((x) => String(x.poolType) === state.activePool)
      if (!p) return 0
      return Math.min(p.pity / 80, 1)
    },
    /** 最近 5★（按激活卡池过滤） */
    recent5Filtered(state): RecentItem[] {
      if (!state.detail) return []
      if (state.activePool === 'all') return state.detail.recent5
      return state.detail.recent5.filter((r) => r.pool === state.activePool)
    },
    /** 激活卡池的出金间隔分布 */
    distFiltered(state): number[] {
      if (!state.detail) return []
      if (state.activePool === 'all') {
        const sum = new Array(9).fill(0)
        for (const arr of Object.values(state.detail.dist)) {
          arr.forEach((v, i) => (sum[i] += v))
        }
        return sum
      }
      return state.detail.dist[state.activePool] ?? new Array(9).fill(0)
    },
    /** 5★ 构成：活动池 vs 常驻类池 */
    star5Split(state): { limited: number; standard: number } {
      if (!state.detail) return { limited: 0, standard: 0 }
      let limited = 0
      let standard = 0
      for (const p of state.detail.stats.pools) {
        if (['1', '2', '8', '9', '10', '11'].includes(String(p.poolType))) limited += p.count5
        else standard += p.count5
      }
      return { limited, standard }
    },
  },

  actions: {
    async init() {
      Events.On('gacha:progress', (ev: any) => {
        this.progress = ev?.data ?? ev
      })
      await this.loadAccounts()
      if (this.accounts.length > 0) {
        await this.openAccount(this.accounts[0].uid)
      }
    },

    async loadAccounts() {
      this.accounts = ((await GachaService.ListAccounts()) as any[]) ?? []
    },

    async openAccount(uid: string) {
      this.detail = ((await GachaService.OpenAccount(uid)) as any) ?? null
      this.activePool = 'all'
    },

    async refresh() {
      this.loading = true
      this.lastError = ''
      this.progress = null
      try {
        this.detail = ((await GachaService.Refresh('', '')) as any) ?? null
        this.activePool = 'all'
        await this.loadAccounts()
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      } finally {
        this.loading = false
        this.progress = null
      }
    },

    async exportAccount() {
      if (!this.detail) return ''
      this.lastError = ''
      try {
        return (await GachaService.ExportAccount(this.detail.uid, '')) as unknown as string
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
        return ''
      }
    },

    switchAccount(uid: string) {
      return this.openAccount(uid)
    },
  },
})
