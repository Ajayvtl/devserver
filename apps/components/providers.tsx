'use client'

import type { ReactNode } from 'react'

import { ToastProvider, ToastViewport } from './toast'

import { AuthProvider } from './auth/auth-context'

interface Props {
  children: ReactNode
}

export function Providers({ children }: Props) {
  return (
    <AuthProvider>
      <ToastProvider>
        {children}
        <ToastViewport />
      </ToastProvider>
    </AuthProvider>
  )
}
