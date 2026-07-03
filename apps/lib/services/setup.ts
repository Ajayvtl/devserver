import { submitCommand } from '../api/client'
import { setupFallback } from './fallbacks'
import type { CommandRecord, SetupWizardData } from '../types'

export async function getSetupWizardData(): Promise<SetupWizardData> {
  return setupFallback
}

export async function completeSetupWizard(payload: {
  licenseAccepted: boolean
  adminName: string
  adminEmail: string
  adminPassword: string
  provider: string
  installationType: string
}): Promise<CommandRecord> {
  return submitCommand<CommandRecord>({
    capability: 'workspace.setup',
    parameters: payload,
    name: 'Bootstrap DevServer',
  }, {
    auth: false,
  })
}
