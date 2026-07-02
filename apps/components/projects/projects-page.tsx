import Link from 'next/link'

import type { ProjectListItem } from '@/lib/types'

import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { DataTable } from '../ui/table'
import { EmptyState } from '../ui/empty-state'
import { MetricCard } from '../ui/metric-card'

interface Props {
  projects: ProjectListItem[]
}

export function ProjectsPage({ projects }: Props) {
  const healthy = projects.filter((project) => project.status === 'healthy').length

  return (
    <div className="page">
      <div className="section-header">
        <div>
          <div className="section-header__eyebrow">Workspaces</div>
          <h1 className="section-header__title">Manage every workspace in one place</h1>
          <p className="section-header__description">
            Create, edit, and inspect workspaces without leaving the product shell.
          </p>
        </div>
        <div className="section-header__actions">
          <Button href="/projects/new" variant="primary">
            New workspace
          </Button>
          <Button href="/settings" variant="secondary">
            Workspace defaults
          </Button>
        </div>
      </div>

      <div className="metric-grid metric-grid--three">
        <MetricCard label="Workspaces" value={`${projects.length}`} detail="All registered workloads" trend="+2" tone="accent" />
        <MetricCard label="Healthy" value={`${healthy}`} detail="Workspaces in good shape" trend="+1" tone="success" />
        <MetricCard label="Needs attention" value={`${projects.length - healthy}`} detail="Warnings and blocked items" trend="stable" tone="warning" />
      </div>

      <Card>
        <div className="card__header">
          <div>
            <div className="card__eyebrow">List</div>
            <h2 className="card__title">Workspaces</h2>
          </div>
          <Badge tone="info">Service-backed</Badge>
        </div>
        {projects.length ? (
          <DataTable
            columns={[
              { header: 'Workspace' },
              { header: 'Owner' },
              { header: 'Environment' },
              { header: 'Status' },
              { header: 'Updated' },
              { header: 'Action' },
            ]}
            rows={projects.map((project) => [
              <div key={project.slug}>
                <strong>{project.name}</strong>
                <p>{project.description}</p>
              </div>,
              project.owner,
              project.environment,
              <Badge key={project.slug} tone={project.status === 'healthy' ? 'success' : project.status === 'blocked' ? 'danger' : 'warning'}>
                {project.status}
              </Badge>,
              project.updatedAt,
              <Link key={project.slug} href={`/projects/${project.slug}`}>
                View
              </Link>,
            ])}
          />
        ) : (
          <EmptyState title="No workspaces yet" description="Create your first workspace to start managing deployments and domains." actionLabel="New workspace" actionHref="/projects/new" />
        )}
      </Card>
    </div>
  )
}
