import { AppShell } from '@/components/app-shell'
import { ProjectForm } from '@/components/projects/project-form'
import { getProjectFormData, getProjectFormOptions } from '@/lib/services/projects'

export default async function NewProjectPage() {
  const [initialData, options] = await Promise.all([getProjectFormData(), getProjectFormOptions('create')])

  return (
    <AppShell server="production-east-1" notifications={4}>
      <ProjectForm initialData={initialData} options={options} />
    </AppShell>
  )
}
