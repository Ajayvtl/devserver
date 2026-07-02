'use client'

import type { ReactNode } from 'react'

import { Button } from './button'

interface Props {
  open: boolean
  title: string
  description: string
  onClose: () => void
  children: ReactNode
}

export function Dialog({ open, title, description, onClose, children }: Props) {
  if (!open) {
    return null
  }

  return (
    <div className="dialog-overlay" role="presentation" onClick={onClose}>
      <div
        className="dialog"
        role="dialog"
        aria-modal="true"
        aria-labelledby="dialog-title"
        aria-describedby="dialog-description"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="dialog__header">
          <div>
            <h3 id="dialog-title">{title}</h3>
            <p id="dialog-description">{description}</p>
          </div>
          <Button variant="ghost" onClick={onClose} type="button">
            Close
          </Button>
        </div>
        <div className="dialog__body">{children}</div>
      </div>
    </div>
  )
}
