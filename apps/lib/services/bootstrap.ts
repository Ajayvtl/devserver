import { requestOrFallback } from '../api/client'
import { bootstrapFallback } from './fallbacks'
import type { BootstrapDecision } from '../types'

export async function getBootstrapDecision(): Promise<BootstrapDecision> {
  return requestOrFallback<BootstrapDecision>('/api/bootstrap/decision', bootstrapFallback)
}
