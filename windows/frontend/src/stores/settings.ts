import { defineStore } from 'pinia'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { SettingsService } = bindings as any

export interface AppSettings {
  autostart: boolean
  signInAuto: boolean
  signInHour: number
  waveNotify: boolean
  logPath: string
  language: string
}

export interface SettingsWithPaths extends AppSettings {
  dataDir: string
  autostartOn: boolean
}

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: null as SettingsWithPaths | null,
    saving: false as boolean,
    lastError: '' as string,
    savedFlash: false as boolean,
  }),

  actions: {
    async load() {
      this.lastError = ''
      try {
        this.settings = ((await SettingsService.Get()) as any) ?? null
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      }
    },

    async save() {
      if (!this.settings) return
      this.saving = true
      this.lastError = ''
      try {
        const s = this.settings
        await SettingsService.Save({
          autostart: s.autostart,
          signInAuto: s.signInAuto,
          signInHour: s.signInHour,
          waveNotify: s.waveNotify,
          logPath: s.logPath,
          language: s.language,
        })
        this.savedFlash = true
        setTimeout(() => (this.savedFlash = false), 1600)
        await this.load()
      } catch (e: any) {
        this.lastError = String(e?.message ?? e)
      } finally {
        this.saving = false
      }
    },
  },
})
