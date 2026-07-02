'use client'

interface Props {
  value: number
  tone?: 'accent' | 'success' | 'warning' | 'info'
}

export function Progress({ value, tone = 'accent' }: Props) {
  return (
    <div className="progress" role="progressbar" aria-valuenow={value} aria-valuemin={0} aria-valuemax={100}>
      <div className={`progress__bar progress__bar--${tone}`} style={{ width: `${value}%` }} />
    </div>
  )
}
