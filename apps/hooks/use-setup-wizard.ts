'use client'

import { useEffect, useState } from 'react'

import type { SetupWizardData } from '@/lib/types'
import { getSetupWizardData } from '@/lib/services/setup'

export function useSetupWizard() {
  const [data, setData] = useState<SetupWizardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    getSetupWizardData().then((next) => {
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
