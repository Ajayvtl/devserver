'use client'

import type { ReactNode } from 'react'

type BadgeTone = 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info'

interface Props {
  tone?: BadgeTone
  children: ReactNode
}

export function Badge({ tone = 'neutral', children }: Props) {
  return <span className={`badge badge-${tone}`}>{children}</span>
}
