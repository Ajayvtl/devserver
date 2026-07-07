'use client'

import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

import { ToastProvider, ToastViewport } from './toast'
import { AuthProvider } from './auth/auth-context'

export type Density = 'comfortable' | 'compact' | 'ultra'
export type FontScale = 14 | 15 | 16

interface PreferenceContextType {
  density: Density
  setDensity: (d: Density) => void
  fontScale: FontScale
  setFontScale: (s: FontScale) => void
}

const PreferenceContext = createContext<PreferenceContextType | undefined>(undefined)

export function usePreferences() {
  const context = useContext(PreferenceContext)
  if (!context) {
    throw new Error('usePreferences must be used within PreferenceProvider')
  }
  return context
}

interface Props {
  children: ReactNode
}

export function Providers({ children }: Props) {
  const [density, setDensity] = useState<Density>('comfortable')
  const [fontScale, setFontScale] = useState<FontScale>(15)

  // Apply density & font scale to document element
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const root = document.documentElement
      
      // Clean previous classes
      root.classList.remove('density-comfortable', 'density-compact', 'density-ultra')
      root.classList.remove('font-scale-14', 'font-scale-15', 'font-scale-16')

      // Apply current settings
      root.classList.add(`density-${density}`)
      root.classList.add(`font-scale-${fontScale}`)
    }
  }, [density, fontScale])

  return (
    <AuthProvider>
      <PreferenceContext.Provider value={{ density, setDensity, fontScale, setFontScale }}>
        <ToastProvider>
          {children}
          <ToastViewport />
        </ToastProvider>
      </PreferenceContext.Provider>
    </AuthProvider>
  )
}
