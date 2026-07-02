'use client'

import { Button } from './button'

interface Props {
  title: string
  description: string
  actionLabel?: string
  actionHref?: string
}

export function EmptyState({ title, description, actionLabel, actionHref }: Props) {
  return (
    <div className="empty-state">
      <h3>{title}</h3>
      <p>{description}</p>
      {actionLabel && actionHref ? <Button href={actionHref}>{actionLabel}</Button> : null}
    </div>
  )
}
