'use client'

import { AppShell } from '@/components/app-shell'
import { SectionHeader } from '@/components/ui/section-header'
import { EnvironmentsPanel } from '@/components/config/environments-panel'

export default function EnvironmentsPage() {
  return (
    <AppShell server="production-east-1" notifications={0}>
      <div className="page">
        <SectionHeader
          eyebrow="Configuration"
          title="Environments & Secrets"
          description="Manage deployment targets, environment variables, and zero-trust secrets."
        />
        <EnvironmentsPanel />
      </div>
    </AppShell>
  )
}
