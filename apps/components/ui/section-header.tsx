import type { ReactNode } from 'react'

interface Props {
  eyebrow?: string
  title: string
  description: string
  actions?: ReactNode
}

export function SectionHeader({ eyebrow, title, description, actions }: Props) {
  return (
    <div className="section-header">
      <div>
        {eyebrow ? <div className="section-header__eyebrow">{eyebrow}</div> : null}
        <h1 className="section-header__title">{title}</h1>
        <p className="section-header__description">{description}</p>
      </div>
      {actions ? <div className="section-header__actions">{actions}</div> : null}
    </div>
  )
}
