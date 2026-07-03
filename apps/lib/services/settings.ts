import { settingsFallback } from './fallbacks'
import type { SettingsData } from '../types'

export async function getSettingsData(): Promise<SettingsData> {
  return settingsFallback
}
