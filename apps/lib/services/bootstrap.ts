import { bootstrapFallback } from './fallbacks'
import type { BootstrapDecision } from '../types'

export async function getBootstrapDecision(): Promise<BootstrapDecision> {
  return bootstrapFallback
}
