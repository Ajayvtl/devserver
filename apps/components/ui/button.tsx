'use client'

import Link from 'next/link'
import type { ReactNode } from 'react'

type Variant = 'primary' | 'secondary' | 'ghost'

interface Props {
  children: ReactNode
  href?: string
  variant?: Variant
  type?: 'button' | 'submit'
  onClick?: () => void
  disabled?: boolean
}

export function Button({ children, href, variant = 'secondary', type = 'button', onClick, disabled }: Props) {
  const className = `button button--${variant}`

  if (href) {
    if (disabled) {
      return (
        <span className={`${className} is-disabled`} aria-disabled="true">
          {children}
        </span>
      )
    }

    if (href.startsWith('#') || href.startsWith('mailto:') || href.startsWith('http')) {
      return (
        <a className={className} href={href}>
          {children}
        </a>
      )
    }

    return (
      <Link className={className} href={href}>
        {children}
      </Link>
    )
  }

  return (
    <button className={className} type={type} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  )
}
