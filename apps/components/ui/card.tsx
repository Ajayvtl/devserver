'use client'

import type { HTMLAttributes, ReactNode } from 'react'

interface Props extends HTMLAttributes<HTMLElement> {
  children: ReactNode
}

export function Card({ children, className = '', ...props }: Props) {
  return (
    <section className={`card ${className}`.trim()} {...props}>
      {children}
    </section>
  )
}
