'use client'

import { useState, useCallback, useEffect } from 'react'
import { request } from '@/lib/api/client'
import { useToast } from '../components/toast'

export interface Member {
  id: string
  userId: string
  username: string
  email: string
  roleId: string
  roleName: string
  status: string
  lastLoginAt?: string
  joinedAt?: string
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: string
  name: string
}

interface UseMembersManagementProps {
  currentOrgId?: string | null
}

export function useMembersManagement({ currentOrgId }: UseMembersManagementProps) {
  const { push } = useToast()

  // Data states
  const [members, setMembers] = useState<Member[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [rolesLoading, setRolesLoading] = useState(true)
  const [rolesError, setRolesError] = useState<string | null>(null)

  // Query parameters
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [perPage] = useState(10)
  const [totalPages, setTotalPages] = useState(1)
  const [totalCount, setTotalCount] = useState(0)
  const [sortField, setSortField] = useState('created_at')
  const [sortDir, setSortDir] = useState<'ASC' | 'DESC'>('DESC')

  // Modals & Actions states
  const [inviteOpen, setInviteOpen] = useState(false)
  const [inviteEmail, setInviteEmail] = useState('')
  const [inviteRoleId, setInviteRoleId] = useState('')
  const [inviting, setInviting] = useState(false)

  const [confirmRemoveOpen, setConfirmRemoveOpen] = useState(false)
  const [targetMember, setTargetMember] = useState<Member | null>(null)
  const [removing, setRemoving] = useState(false)

  const [confirmTransferOpen, setConfirmTransferOpen] = useState(false)
  const [transferring, setTransferring] = useState(false)

  const loadRoles = useCallback(async () => {
    if (!currentOrgId) return
    setRolesLoading(true)
    setRolesError(null)
    try {
      const response = await request<Role[]>('/api/v1/roles', {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setRoles(response || [])
      if (response && response.length > 0 && !inviteRoleId) {
        setInviteRoleId(response[0].id)
      }
    } catch (err: any) {
      setRolesError(err.message || 'Failed to load organization roles.')
    } finally {
      setRolesLoading(false)
    }
  }, [currentOrgId, inviteRoleId])

  const loadMembers = useCallback(async () => {
    if (!currentOrgId) return
    setLoading(true)
    setError(null)
    try {
      const query = `/api/v1/organizations/members?page=${page}&per_page=${perPage}&sort=${sortField}&dir=${sortDir}&query=${encodeURIComponent(search)}`
      const response = await request<any>(query, {
        headers: { 'X-Org-ID': currentOrgId }
      })
      if (response && response.data) {
        setMembers(response.data)
        if (response.meta) {
          setTotalPages(response.meta.totalPages || 1)
          setTotalCount(response.meta.total || response.data.length)
        }
      } else {
        setMembers([])
        setTotalPages(1)
        setTotalCount(0)
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load organization members.')
    } finally {
      setLoading(false)
    }
  }, [currentOrgId, page, perPage, sortField, sortDir, search])

  // Loaders Trigger
  useEffect(() => {
    loadRoles()
  }, [loadRoles])

  useEffect(() => {
    loadMembers()
  }, [loadMembers])

  const handleInvite = useCallback(async (e?: React.FormEvent) => {
    if (e) e.preventDefault()
    if (!inviteEmail.trim() || !inviteRoleId || !currentOrgId) return
    setInviting(true)
    try {
      await request('/api/v1/organizations/members', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { email: inviteEmail.trim(), roleId: inviteRoleId }
      })
      push({ title: 'Success', message: `Invitation sent to ${inviteEmail}`, tone: 'success' })
      setInviteOpen(false)
      setInviteEmail('')
      loadMembers()
    } catch (err: any) {
      push({ title: 'Invitation Failed', message: err.message, tone: 'danger' })
    } finally {
      setInviting(false)
    }
  }, [inviteEmail, inviteRoleId, currentOrgId, push, loadMembers])

  const handleRemove = useCallback(async () => {
    if (!targetMember || !currentOrgId) return
    setRemoving(true)
    try {
      await request(`/api/v1/organizations/members?userId=${targetMember.userId}`, {
        method: 'DELETE',
        headers: { 'X-Org-ID': currentOrgId }
      })
      push({ title: 'Removed', message: `${targetMember.username || targetMember.email} has been removed.`, tone: 'success' })
      setConfirmRemoveOpen(false)
      setTargetMember(null)
      loadMembers()
    } catch (err: any) {
      push({ title: 'Removal Failed', message: err.message, tone: 'danger' })
    } finally {
      setRemoving(false)
    }
  }, [targetMember, currentOrgId, push, loadMembers])

  const handleRoleChange = useCallback(async (member: Member, roleId: string) => {
    if (!currentOrgId) return
    try {
      await request('/api/v1/organizations/members', {
        method: 'PUT',
        headers: { 'X-Org-ID': currentOrgId },
        body: { userId: member.userId, roleId }
      })
      push({ title: 'Role Updated', message: 'Member role updated successfully.', tone: 'success' })
      loadMembers()
    } catch (err: any) {
      push({ title: 'Update Failed', message: err.message, tone: 'danger' })
    }
  }, [currentOrgId, push, loadMembers])

  const handleStatusChange = useCallback(async (member: Member, status: string) => {
    if (!currentOrgId) return
    try {
      await request('/api/v1/organizations/members', {
        method: 'PUT',
        headers: { 'X-Org-ID': currentOrgId },
        body: { userId: member.userId, status }
      })
      push({ title: 'Status Updated', message: `Member status set to ${status}.`, tone: 'success' })
      loadMembers()
    } catch (err: any) {
      push({ title: 'Update Failed', message: err.message, tone: 'danger' })
    }
  }, [currentOrgId, push, loadMembers])

  const handleTransferOwnership = useCallback(async () => {
    if (!targetMember || !currentOrgId) return
    setTransferring(true)
    try {
      await request('/api/v1/organizations/transfer-ownership', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { newOwnerUserId: targetMember.userId }
      })
      push({ title: 'Ownership Transferred', message: 'You have transferred ownership of this organization.', tone: 'success' })
      setConfirmTransferOpen(false)
      setTargetMember(null)
      loadMembers()
    } catch (err: any) {
      push({ title: 'Transfer Failed', message: err.message, tone: 'danger' })
    } finally {
      setTransferring(false)
    }
  }, [targetMember, currentOrgId, push, loadMembers])

  const handleSort = useCallback((field: string) => {
    setSortField((prevField) => {
      if (prevField === field) {
        setSortDir((prevDir) => (prevDir === 'ASC' ? 'DESC' : 'ASC'))
      } else {
        setSortDir('ASC')
      }
      return field
    })
    setPage(1)
  }, [])

  return {
    members,
    roles,
    loading,
    error,
    rolesLoading,
    rolesError,
    search,
    setSearch,
    page,
    setPage,
    perPage,
    totalPages,
    totalCount,
    sortField,
    sortDir,
    inviteOpen,
    setInviteOpen,
    inviteEmail,
    setInviteEmail,
    inviteRoleId,
    setInviteRoleId,
    inviting,
    confirmRemoveOpen,
    setConfirmRemoveOpen,
    confirmTransferOpen,
    setConfirmTransferOpen,
    targetMember,
    setTargetMember,
    removing,
    transferring,
    loadRoles,
    loadMembers,
    handleInvite,
    handleRemove,
    handleRoleChange,
    handleStatusChange,
    handleTransferOwnership,
    handleSort
  }
}
