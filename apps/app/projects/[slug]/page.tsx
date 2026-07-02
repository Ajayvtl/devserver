import { notFound } from 'next/navigation'

import { AppShell } from '@/components/app-shell'
import { ProjectDetailView } from '@/components/projects/project-detail-view'
import { getProjectDetail } from '@/lib/services/projects'

interface Props {
  params: Promise<{ slug: string }>
}

export default async function ProjectDetailPage({ params }: Props) {
  const { slug } = await params
  const data = await getProjectDetail(slug)

  if (!data.project) {
    notFound()
  }

  return (
    <AppShell server="production-east-1" notifications={4}>
      <ProjectDetailView data={data} />
    </AppShell>
  )
}
