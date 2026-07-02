import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { CodeBlock } from '../ui/code-block'
import { EmptyState } from '../ui/empty-state'
import { LearnCard } from '../ui/learn-card'
import { SectionHeader } from '../ui/section-header'
import type { DevCenterSectionData } from '@/lib/types'

interface Props {
  data: DevCenterSectionData
}

export function DevCenterPage({ data }: Props) {
  return (
    <div className="page">
      <SectionHeader
        eyebrow="DevCenter"
        title={data.title}
        description={data.subtitle}
        actions={<Badge tone="accent">{data.metadata?.project ? `Project: ${data.metadata.project}` : 'Live content'}</Badge>}
      />

      <div className="page-grid--two">
        <Card>
          <div className="card__eyebrow">Why it matters</div>
          <p className="page-subtitle">{data.description}</p>
          <ul className="bullet-list">
            {data.bullets.map((bullet) => (
              <li key={bullet}>{bullet}</li>
            ))}
          </ul>
        </Card>
        <Card>
          <div className="card__eyebrow">Examples</div>
          <div className="stack">
            {data.examples.map((example) => (
              <div key={example} className="detail-item">
                <strong>{example}</strong>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <div className="page-grid--two">
        <Card>
          <div className="card__eyebrow">Checklist</div>
          <div className="check-list">
            {data.checklist.map((item) => (
              <div key={item.label} className="check-item">
                <div>
                  <div className="check-item__label">
                    <strong>{item.label}</strong>
                    <Badge tone={toneForStatus(item.status)}>{item.status}</Badge>
                  </div>
                  <p className="check-item__detail">{item.detail}</p>
                </div>
              </div>
            ))}
          </div>
        </Card>
        <Card>
          <div className="card__eyebrow">Code sample</div>
          <CodeBlock title={data.title}>{data.codeSample}</CodeBlock>
        </Card>
      </div>

      <div className="page-grid--two">
        {data.knowledge.length ? <LearnCard title={data.title} articles={data.knowledge} /> : <EmptyState title={`${data.title} placeholder`} description="This devcenter section is ready to be expanded with real docs." />}
        <Card>
          <div className="card__eyebrow">Section key</div>
          <h3>{data.key}</h3>
          <p>{data.subtitle}</p>
        </Card>
      </div>

      {data.metadata && Object.keys(data.metadata).length ? (
        <div className="page-grid--two">
          <Card>
            <div className="card__eyebrow">Project metadata</div>
            <div className="stack">
              {Object.entries(data.metadata).map(([key, value]) => (
                <div key={key} className="detail-item">
                  <span>{key}</span>
                  <strong>{value}</strong>
                </div>
              ))}
            </div>
          </Card>
        </div>
      ) : null}
    </div>
  )
}

function toneForStatus(status: string) {
  if (status === 'passed') return 'success'
  if (status === 'warning') return 'warning'
  if (status === 'failed') return 'danger'
  return 'info'
}
