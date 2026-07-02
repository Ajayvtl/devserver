import { Badge } from './badge'

interface Props {
  label: string
  value: string
  detail: string
  trend: string
  tone: 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info'
}

export function MetricCard({ label, value, detail, trend, tone }: Props) {
  return (
    <div className={`metric-card metric-${tone}`}>
      <div className="metric-card__header">
        <span className="metric-card__label">{label}</span>
        <Badge tone={tone}>{trend}</Badge>
      </div>
      <div className="metric-card__value">{value}</div>
      <p className="metric-card__helper">{detail}</p>
    </div>
  )
}
