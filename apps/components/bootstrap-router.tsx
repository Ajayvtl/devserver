'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

import { useBootstrapDecision } from '@/hooks/use-bootstrap-decision'

import { BootstrapScreen } from './bootstrap-screen'

export function BootstrapRouter() {
  const router = useRouter()
  const { decision, loading } = useBootstrapDecision()

  useEffect(() => {
    if (!loading && decision) {
      const timer = window.setTimeout(() => {
        router.replace(`/${decision.destination}`)
      }, 1400)

      return () => window.clearTimeout(timer)
    }
  }, [decision, loading, router])

  return (
    <BootstrapScreen
      progress={loading ? 32 : Math.max(60, decision?.installationProgress ?? 0)}
      message={decision?.message ?? 'Checking installation prerequisites and state store...'}
    />
  )
}
