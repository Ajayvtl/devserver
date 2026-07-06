'use client'

import { AppShell } from '@/components/app-shell'
import { SectionHeader } from '@/components/ui/section-header'
import { ProvidersPanel } from '@/components/config/providers-panel'

export default function ProvidersPage() {
  return (
    <AppShell server="production-east-1" notifications={0}>
      <div className="page">
        <SectionHeader
          eyebrow="Configuration"
          title="AI Providers"
          description="Manage API keys, endpoints, and default models for LLM inference engines."
        />
        <ProvidersPanel />
      </div>
    </AppShell>
  )
}
