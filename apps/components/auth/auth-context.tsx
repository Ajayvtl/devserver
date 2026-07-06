'use client'

import React, { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { getApiToken, request } from '@/lib/api/client'
import { logout } from '@/lib/services/auth'

interface UserProfile {
  id: string
  status: string
}

interface Organization {
  id: string
  name: string
}

interface AuthContextType {
  user: UserProfile | null
  organizations: Organization[]
  currentOrgId: string | null
  setCurrentOrgId: (id: string) => void
  isLoading: boolean
  signOut: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter()
  const pathname = usePathname()
  const [user, setUser] = useState<UserProfile | null>(null)
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [currentOrgId, setCurrentOrgId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    async function loadAuth() {
      const token = getApiToken()
      if (!token) {
        setIsLoading(false)
        if (pathname !== '/login' && pathname !== '/setup') {
          router.replace('/login')
        }
        return
      }

      try {
        const userProfile = await request<UserProfile>('/api/v1/users/me')
        setUser(userProfile)

        const orgs = await request<Organization[]>('/api/v1/organizations')
        setOrganizations(orgs)

        if (orgs.length > 0) {
          const storedOrg = typeof window !== 'undefined' ? window.localStorage.getItem('devserver-org') : null
          if (storedOrg && orgs.find(o => o.id === storedOrg)) {
            setCurrentOrgId(storedOrg)
          } else {
            setCurrentOrgId(orgs[0].id)
            if (typeof window !== 'undefined') window.localStorage.setItem('devserver-org', orgs[0].id)
          }
        }
      } catch (error) {
        console.error('Auth error', error)
        setUser(null)
        logout()
        if (pathname !== '/login' && pathname !== '/setup') {
          router.replace('/login')
        }
      } finally {
        setIsLoading(false)
      }
    }

    loadAuth()
  }, [pathname, router])

  const handleSetOrgId = (id: string) => {
    setCurrentOrgId(id)
    if (typeof window !== 'undefined') {
      window.localStorage.setItem('devserver-org', id)
    }
  }

  const signOut = async () => {
    await logout()
    setUser(null)
    setOrganizations([])
    setCurrentOrgId(null)
    router.replace('/login')
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        organizations,
        currentOrgId,
        setCurrentOrgId: handleSetOrgId,
        isLoading,
        signOut
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
