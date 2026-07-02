import { request, requestOrFallback } from '../api/client'
import { setupFallback } from './fallbacks'
import type { SetupWizardData } from '../types'

export async function getSetupWizardData(): Promise<SetupWizardData> {
  return requestOrFallback<SetupWizardData>('/api/setup', setupFallback)
}

export async function completeSetupWizard(payload: {
  licenseAccepted: boolean
  adminName: string
  adminEmail: string
  adminPassword: string
  provider: string
  installationType: string
}): Promise<{ taskId: string; status: string }> {
  return request<{ taskId: string; status: string }>('/api/setup/complete', {
    method: 'POST',
    body: payload,
    auth: false,
  })
}
