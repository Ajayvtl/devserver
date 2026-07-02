'use client'

import { useEffect, useState } from 'react'

import type { LoginData } from '@/lib/types'
import { getLoginData } from '@/lib/services/auth'

export function useLoginData() {
  const [data, setData] = useState<LoginData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true

    getLoginData().then((next) => {
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
