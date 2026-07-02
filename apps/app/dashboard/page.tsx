import { getDashboardData } from '@/lib/services/dashboard'
import { getKnowledgeArticles } from '@/lib/services/knowledge'

import { AppShell } from '@/components/app-shell'
import { DashboardPageContent } from '@/components/dashboard/dashboard-page'

export default async function DashboardPage() {
  const [data, knowledge] = await Promise.all([getDashboardData(), getKnowledgeArticles('dashboard')])

  return (
    <AppShell server={data.server} notifications={data.notifications}>
      <DashboardPageContent data={data} knowledge={knowledge} />
    </AppShell>
  )
}
