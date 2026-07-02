'use client'

import { useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'

import type { ProjectFormData, ProjectFormOptions } from '@/lib/types'

import { saveProject } from '@/lib/services/projects'

import { useToast } from '../toast'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { Input } from '../ui/input'
import { Textarea } from '../ui/textarea'

interface Props {
  initialData: ProjectFormData
  options: ProjectFormOptions
}

export function ProjectForm({ initialData, options }: Props) {
  const router = useRouter()
  const { push } = useToast()
  const [form, setForm] = useState(initialData)
  const [errors, setErrors] = useState<string[]>([])
  const [saving, setSaving] = useState(false)

  const canSubmit = useMemo(() => {
    return (
      form.name.trim().length > 1 &&
      form.slug.trim().length > 2 &&
      form.repository.includes('git') &&
      form.owner.trim().length > 1 &&
      form.domains.trim().length > 0
    )
  }, [form])

  const validate = () => {
    const next: string[] = []
    if (form.name.trim().length < 2) next.push('Project name is required.')
    if (form.slug.trim().length < 3) next.push('Project slug must be at least 3 characters.')
    if (!form.repository.includes('git')) next.push('Repository should be a valid git URL.')
    if (form.owner.trim().length < 2) next.push('Owner is required.')
    if (!form.domains.trim()) next.push('At least one domain should be defined.')
    setErrors(next)
    return next.length === 0
  }

  const submit = async () => {
    if (!validate()) {
      push({
        title: 'Project form needs attention',
        message: 'Please fix the highlighted fields before continuing.',
        tone: 'warning',
      })
      return
    }

    setSaving(true)
    const mode = options.submitLabel === 'Save project' ? 'edit' : 'create'
    const saved = await saveProject(form, mode)
    setSaving(false)
    push({
      title: options.title,
      message: `${saved.name} is ready in the mock service.`,
      tone: 'success',
    })
    router.push(`/projects/${saved.slug}`)
  }

  return (
    <div className="form-layout">
      <div className="form-layout__main">
        <Card>
          <div className="card__eyebrow">Project management</div>
          <h1 className="page-title">{options.title}</h1>
          <p className="page-subtitle">{options.subtitle}</p>
          <div className="form-badges">
            <Badge tone="accent">Mock service</Badge>
            <Badge tone="info">Accessible form</Badge>
            <Badge tone="success">Responsive</Badge>
          </div>
        </Card>

        {errors.length ? (
          <Card className="wizard-errors">
            {errors.map((error) => (
              <div key={error} className="wizard-error">
                {error}
              </div>
            ))}
          </Card>
        ) : null}

        <Card>
          <div className="form-grid">
            <Input label="Project name" value={form.name} onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))} />
            <Input label="Slug" value={form.slug} onChange={(event) => setForm((current) => ({ ...current, slug: event.target.value }))} />
            <Input label="Owner" value={form.owner} onChange={(event) => setForm((current) => ({ ...current, owner: event.target.value }))} />
            <Input label="Environment" value={form.environment} onChange={(event) => setForm((current) => ({ ...current, environment: event.target.value }))} />
            <Input label="Repository" value={form.repository} onChange={(event) => setForm((current) => ({ ...current, repository: event.target.value }))} />
            <Input label="Branch" value={form.branch} onChange={(event) => setForm((current) => ({ ...current, branch: event.target.value }))} />
            <Input label="Deployment target" value={form.deployTarget} onChange={(event) => setForm((current) => ({ ...current, deployTarget: event.target.value }))} />
            <Input label="Health check path" value={form.healthCheck} onChange={(event) => setForm((current) => ({ ...current, healthCheck: event.target.value }))} />
          </div>
          <Textarea
            label="Description"
            rows={5}
            value={form.description}
            onChange={(event) => setForm((current) => ({ ...current, description: event.target.value }))}
          />
          <Input
            label="Domains"
            value={form.domains}
            onChange={(event) => setForm((current) => ({ ...current, domains: event.target.value }))}
            placeholder="app.example.com, api.example.com"
          />
          <label className="checkbox-card">
            <input
              type="checkbox"
              checked={form.autoDeploy}
              onChange={(event) => setForm((current) => ({ ...current, autoDeploy: event.target.checked }))}
            />
            <span>Enable automatic deployment for approved commits.</span>
          </label>

          <div className="page-actions">
            <Button href="/projects" variant="ghost">
              Cancel
            </Button>
            <Button type="button" variant="primary" disabled={saving || !canSubmit} onClick={submit}>
              {saving ? 'Saving...' : options.submitLabel}
            </Button>
          </div>
        </Card>
      </div>

      <aside className="form-layout__rail">
        <Card>
          <div className="card__eyebrow">Preview</div>
          <h3>{form.name || 'Untitled project'}</h3>
          <p>{form.description || 'Project description will appear here.'}</p>
          <div className="rail-summary">
            <div>
              <span>Slug</span>
              <strong>{form.slug || 'project-slug'}</strong>
            </div>
            <div>
              <span>Environment</span>
              <strong>{form.environment}</strong>
            </div>
            <div>
              <span>Owner</span>
              <strong>{form.owner || 'Unassigned'}</strong>
            </div>
          </div>
        </Card>
      </aside>
    </div>
  )
}
