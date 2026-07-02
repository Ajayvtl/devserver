'use client'

import type { InputHTMLAttributes } from 'react'

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export function Input({ label, className = '', ...props }: Props) {
  return (
    <label className="field">
      {label ? <span className="field__label">{label}</span> : null}
      <input className={`input ${className}`.trim()} {...props} />
    </label>
  )
}
