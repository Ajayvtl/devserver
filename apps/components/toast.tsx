'use client'

import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'

type ToastTone = 'info' | 'success' | 'warning' | 'danger'

interface ToastItem {
  id: string
  title: string
  message: string
  tone: ToastTone
}

interface ToastContextValue {
  push: (toast: Omit<ToastItem, 'id'>) => void
}

const ToastContext = createContext<ToastContextValue | null>(null)

export function ToastProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([])

  const push = (toast: Omit<ToastItem, 'id'>) => {
    const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`
    setItems((current) => [...current, { ...toast, id }])
  }

  const value = useMemo(() => ({ push }), [])

  useEffect(() => {
    if (items.length === 0) {
      return
    }

    const timeout = window.setTimeout(() => {
      setItems((current) => current.slice(1))
    }, 3200)

    return () => window.clearTimeout(timeout)
  }, [items])

  return <ToastContext.Provider value={value}>{children}{items.length ? <ToastStack items={items} /> : null}</ToastContext.Provider>
}

export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}

function ToastStack({ items }: { items: ToastItem[] }) {
  return (
    <div className="toast-stack" aria-live="polite" aria-atomic="true">
      {items.map((item) => (
        <div key={item.id} className={`toast toast--${item.tone}`}>
          <strong>{item.title}</strong>
          <p>{item.message}</p>
        </div>
      ))}
    </div>
  )
}

export function ToastViewport() {
  return null
}
