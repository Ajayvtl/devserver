'use client'

import type { ReactNode } from 'react'

import { ToastProvider, ToastViewport } from './toast'

interface Props {
  children: ReactNode
}

export function Providers({ children }: Props) {
  return (
    <ToastProvider>
      {children}
      <ToastViewport />
    </ToastProvider>
  )
}
