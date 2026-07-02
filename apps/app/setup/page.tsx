import { getSetupWizardData } from '@/lib/services/setup'

import { SetupWizard } from '@/components/setup/setup-wizard'

export default async function SetupPage() {
  const data = await getSetupWizardData()
  return <SetupWizard initialData={data} />
}
