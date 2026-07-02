'use client'

import type { TextareaHTMLAttributes } from 'react'

interface Props extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string
}

export function Textarea({ label, className = '', ...props }: Props) {
  return (
    <label className="field">
      {label ? <span className="field__label">{label}</span> : null}
      <textarea className={`input textarea ${className}`.trim()} {...props} />
    </label>
  )
}
