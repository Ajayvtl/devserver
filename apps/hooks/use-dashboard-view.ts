'use client'

import { useEffect, useState } from 'react'

import type { DashboardDataV2 } from '@/lib/types'
import { getDashboardData } from '@/lib/services/dashboard'

export function useDashboardView() {
  const [data, setData] = useState<DashboardDataV2 | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    getDashboardData().then((next) => {
      if (!mounted) {
        return
      }
      setData(next)
      setLoading(false)
    })

    return () => {
      mounted = false
    }
  }, [])

  return { data, loading }
}
