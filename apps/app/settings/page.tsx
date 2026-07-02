import { AppShell } from '@/components/app-shell'
import { SettingsPage } from '@/components/settings/settings-page'
import { getSettingsData } from '@/lib/services/settings'

export default async function SettingsIndexPage() {
  const data = await getSettingsData()

  return (
    <AppShell server="production-east-1" notifications={4}>
      <SettingsPage data={data} />
    </AppShell>
  )
}
