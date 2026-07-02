'use client'

import { Button } from './button'

interface Props {
  title: string
  description: string
  actionLabel?: string
  onAction?: () => void
}

export function ErrorState({ title, description, actionLabel, onAction }: Props) {
  return (
    <div className="error-state">
      <h3>{title}</h3>
      <p>{description}</p>
      {actionLabel ? <Button variant="primary" onClick={onAction}>{actionLabel}</Button> : null}
    </div>
  )
}
