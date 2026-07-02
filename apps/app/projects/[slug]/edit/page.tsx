import { AppShell } from '@/components/app-shell'
import { ProjectForm } from '@/components/projects/project-form'
import { getProjectFormData, getProjectFormOptions } from '@/lib/services/projects'

interface Props {
  params: Promise<{ slug: string }>
}

export default async function EditProjectPage({ params }: Props) {
  const { slug } = await params
  const [initialData, options] = await Promise.all([getProjectFormData(slug), getProjectFormOptions('edit')])

  return (
    <AppShell server="production-east-1" notifications={4}>
      <ProjectForm initialData={initialData} options={options} />
    </AppShell>
  )
}
