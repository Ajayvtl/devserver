import { notFound } from 'next/navigation'

import { AppShell } from '@/components/app-shell'
import { DevCenterPage } from '@/components/devcenter/devcenter-page'
import { getDevCenterSection } from '@/lib/services/devcenter'

interface Props {
  params: Promise<{ section: string }>
}

export default async function DevCenterSectionPage({ params }: Props) {
  const { section } = await params
  const data = await getDevCenterSection(section)

  if (!data) {
    notFound()
  }

  return (
    <AppShell server="production-east-1" notifications={4}>
      <DevCenterPage data={data} />
    </AppShell>
  )
}
