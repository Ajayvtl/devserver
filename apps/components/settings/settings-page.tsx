import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { DataTable } from '../ui/table'
import { SectionHeader } from '../ui/section-header'
import type { SettingsData } from '@/lib/types'

interface Props {
  data: SettingsData
}

export function SettingsPage({ data }: Props) {
  return (
    <div className="page">
      <SectionHeader
        eyebrow="Settings"
        title="Global settings"
        description="Control platform preferences, security posture, and notifications from one place."
        actions={<Button variant="primary">Save changes</Button>}
      />

      <div className="page-grid--two">
        {data.sections.map((section) => (
          <Card key={section.title}>
            <div className="card__eyebrow">{section.title}</div>
            <p className="page-subtitle">{section.description}</p>
            <div className="stack">
              {section.items.map((item) => (
                <div key={item.label} className="detail-item">
                  <span>{item.label}</span>
                  <strong>{item.value}</strong>
                  <p>{item.detail}</p>
                </div>
              ))}
            </div>
          </Card>
        ))}
      </div>

      <div className="page-grid--two">
        <Card>
          <div className="card__eyebrow">Preferences</div>
          <DataTable
            columns={[
              { header: 'Preference' },
              { header: 'Value' },
              { header: 'Detail' },
            ]}
            rows={data.preferences.map((item) => [item.label, <Badge key={item.label} tone="accent">{item.value}</Badge>, item.detail])}
          />
        </Card>
        <Card>
          <div className="card__eyebrow">Security</div>
          <DataTable
            columns={[
              { header: 'Control' },
              { header: 'Status' },
              { header: 'Detail' },
            ]}
            rows={data.security.map((item) => [item.label, <Badge key={item.label} tone="info">{item.value}</Badge>, item.detail])}
          />
        </Card>
      </div>

      <div className="page-grid--two">
        <Card>
          <div className="card__eyebrow">AI providers</div>
          <DataTable
            columns={[
              { header: 'Provider' },
              { header: 'Status' },
              { header: 'Detail' },
            ]}
            rows={(data.aiProviders ?? []).map((item) => [
              item.label,
              <Badge key={item.label} tone={item.value === 'Connected' ? 'success' : item.value === 'Disabled' ? 'warning' : 'info'}>
                {item.value}
              </Badge>,
              item.detail,
            ])}
          />
        </Card>
        <Card>
          <div className="card__eyebrow">MCP integrations</div>
          <DataTable
            columns={[
              { header: 'Integration' },
              { header: 'Status' },
              { header: 'Detail' },
            ]}
            rows={(data.mcpIntegrations ?? []).map((item) => [
              item.label,
              <Badge key={item.label} tone={item.value === 'Connected' ? 'success' : 'info'}>
                {item.value}
              </Badge>,
              item.detail,
            ])}
          />
        </Card>
      </div>
    </div>
  )
}
