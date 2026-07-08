'use client'

import { AppShell } from '@/components/app-shell'
import { SuperAdminOrgsList } from '@/components/superadmin/orgs-list'

export default function SuperAdminOrgsPage() {
  return (
    <AppShell server="production-east-1" notifications={1}>
      <SuperAdminOrgsList />
    </AppShell>
  )
}
