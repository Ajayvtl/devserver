import { AppShell } from '@/components/app-shell'
import { SettingsPage } from '@/components/settings/settings-page'

export default function SettingsIndexPage() {
  return (
    <AppShell server="production-east-1" notifications={0}>
      <SettingsPage />
    </AppShell>
  )
}
