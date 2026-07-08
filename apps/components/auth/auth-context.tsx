'use client'

import React, { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { getApiToken, request } from '@/lib/api/client'
import { logout } from '@/lib/services/auth'

interface UserProfile {
  id: string
  status: string
  role?: string
  permissions?: string[]
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
  can: (permission: string) => boolean
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
        const storedOrg = typeof window !== 'undefined' ? window.localStorage.getItem('devserver-org') : null

        const userProfile = await request<UserProfile>('/api/v1/users/me', {
          headers: storedOrg ? { 'X-Org-ID': storedOrg } : undefined
        })
        setUser(userProfile)

        const orgsResp = await request<Organization[]>('/api/v1/organizations')
        const orgs = orgsResp || []
        setOrganizations(orgs)

        if (orgs.length > 0) {
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

    const handleSessionExpired = () => {
      setUser(null)
      setOrganizations([])
      setCurrentOrgId(null)
      logout()

      // Push friendly error notification
      if (typeof window !== 'undefined') {
        const nextUrl = window.location.pathname
        router.replace(`/login?redirect=${encodeURIComponent(nextUrl)}&expired=true`)
      }
    }

    if (typeof window !== 'undefined') {
      window.addEventListener('devserver-session-expired', handleSessionExpired)
    }

    return () => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('devserver-session-expired', handleSessionExpired)
      }
    }
  }, [pathname, router])

  // Update role profile dynamically when organization changes
  useEffect(() => {
    if (!currentOrgId || !user) return
    async function updateRole() {
      try {
        const userProfile = await request<UserProfile>('/api/v1/users/me', {
          headers: { 'X-Org-ID': currentOrgId as string }
        })
        setUser(prev => prev ? { ...prev, role: userProfile.role, permissions: userProfile.permissions } : userProfile)
      } catch (err) {
        console.error('Failed to update user role', err)
      }
    }
    updateRole()
  }, [currentOrgId, user])

  const handleSetOrgId = (id: string) => {
    setCurrentOrgId(id)
    if (typeof window !== 'undefined') {
      window.localStorage.setItem('devserver-org', id)
      // Invalidate caches and trigger re-render of tenant data
      window.dispatchEvent(new Event('devserver-org-changed'))
    }
  }

  const can = (permission: string) => {
    if (!user || !user.permissions) return false
    return user.permissions.includes(permission) || user.permissions.includes('*')
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
        signOut,
        can
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
