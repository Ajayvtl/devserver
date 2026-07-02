import { requestOrFallback } from '../api/client'
import { settingsFallback } from './fallbacks'
import type { SettingsData } from '../types'

export async function getSettingsData(): Promise<SettingsData> {
  return requestOrFallback<SettingsData>('/api/settings', settingsFallback)
}
