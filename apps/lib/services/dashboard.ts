import { requestOrFallback } from '../api/client'
import { dashboardFallback } from './fallbacks'
import type { DashboardDataV2 } from '../types'

export async function getDashboardData(): Promise<DashboardDataV2> {
  return requestOrFallback<DashboardDataV2>('/api/dashboard', dashboardFallback)
}
