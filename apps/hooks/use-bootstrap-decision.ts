'use client'

import { useEffect, useState } from 'react'

import type { BootstrapDecision } from '@/lib/types'
import { getBootstrapDecision } from '@/lib/services/bootstrap'

export function useBootstrapDecision() {
  const [decision, setDecision] = useState<BootstrapDecision | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    getBootstrapDecision().then((result) => {
      if (!mounted) {
        return
      }
      setDecision(result)
      setLoading(false)
    })

    return () => {
      mounted = false
    }
  }, [])

  return { decision, loading }
}
