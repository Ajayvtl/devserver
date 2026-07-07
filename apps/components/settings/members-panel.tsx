'use client'

import { useAuth } from '../auth/auth-context'
import { Card } from '../ui/card'
import { Button } from '../ui/button'
import { DataTable } from '../ui/table'
import { Input } from '../ui/input'
import { Badge } from '../ui/badge'
import { Dialog } from '../ui/dialog'
import { Skeleton } from '../ui/skeleton'
import { EmptyState } from '../ui/empty-state'
import { ErrorState } from '../ui/error-state'
import { useMembersManagement } from '../../hooks/use-members-management'

export function MembersPanel() {
  const { currentOrgId, can, user } = useAuth()

  const {
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
  } = useMembersManagement({ currentOrgId })

  const formatTimestamp = (ts?: string) => {
    if (!ts) return 'Never'
    return new Date(ts).toLocaleString(undefined, {
      dateStyle: 'short',
      timeStyle: 'short'
    })
  }

  const getStatusTone = (status: string) => {
    switch (status) {
      case 'active':
        return 'success'
      case 'invited':
        return 'info'
      case 'suspended':
      case 'deactivated':
        return 'danger'
      default:
        return 'neutral'
    }
  }

  const isOwner = can('org.owner')
  const hasInviteAccess = can('members.invite')
  const hasManageAccess = can('members.manage')
  const hasRemoveAccess = can('members.remove')

  // Error & loading views for Roles load state
  if (rolesLoading) {
    return <Skeleton lines={5} />
  }

  if (rolesError) {
    return (
      <ErrorState
        title="Failed to Load Roles"
        description={rolesError}
        actionLabel="Retry"
        onAction={loadRoles}
      />
    )
  }

  return (
    <div className="stack">
      <Card>
        <div className="card__header">
          <div>
            <h2 className="card__title">Organization Members</h2>
            <p className="card__eyebrow">
              Manage access, modify user roles, and monitor team activity.
            </p>
          </div>
          {hasInviteAccess && (
            <Button variant="primary" onClick={() => setInviteOpen(true)}>
              Invite Member
            </Button>
          )}
        </div>

        {/* Search and Sort controls */}
        <div className="page-actions">
          <Input
            placeholder="Search by username or email..."
            value={search}
            onChange={e => {
              setSearch(e.target.value)
              setPage(1)
            }}
          />

          <div className="page-actions">
            <label className="field">
              <span className="field__label">Sort By</span>
              <select
                className="select"
                value={sortField}
                onChange={e => handleSort(e.target.value)}
              >
                <option value="username">Username</option>
                <option value="email">Email</option>
                <option value="role_id">Role</option>
                <option value="status">Status</option>
                <option value="created_at">Joined Date</option>
              </select>
            </label>

            <label className="field">
              <span className="field__label">Order</span>
              <select
                className="select"
                value={sortDir}
                onChange={e => handleSort(sortField)}
              >
                <option value="ASC">Ascending</option>
                <option value="DESC">Descending</option>
              </select>
            </label>
          </div>
        </div>

        {loading ? (
          <Skeleton lines={5} />
        ) : error ? (
          <ErrorState
            title="Error Loading Members"
            description={error}
            actionLabel="Retry"
            onAction={loadMembers}
          />
        ) : members.length === 0 ? (
          <EmptyState
            title="No Members Found"
            description="No users matched your search criteria or the organization has no other members."
          />
        ) : (
          <div className="stack">
            <DataTable
              columns={[
                { header: 'Username' },
                { header: 'Email' },
                { header: 'Role' },
                { header: 'Status' },
                { header: 'Last Login' },
                { header: 'Joined' },
                { header: 'Actions' }
              ]}
              rows={members.map(member => [
                <span key={`username-${member.id}`} className="sidebar__link-label">
                  {member.username}
                </span>,
                <span key={`email-${member.id}`}>{member.email}</span>,
                hasManageAccess && member.userId !== user?.id ? (
                  <select
                    key={`role-${member.id}`}
                    value={member.roleId}
                    onChange={e => handleRoleChange(member, e.target.value)}
                    className="select"
                  >
                    {roles.map(r => (
                      <option key={r.id} value={r.id}>
                        {r.name}
                      </option>
                    ))}
                  </select>
                ) : (
                  <span key={`role-label-${member.id}`} className="card__eyebrow">
                    {member.roleName || 'Member'}
                  </span>
                ),
                <Badge key={`status-${member.id}`} tone={getStatusTone(member.status)}>
                  {member.status}
                </Badge>,
                <span key={`login-${member.id}`} className="card__eyebrow">
                  {formatTimestamp(member.lastLoginAt)}
                </span>,
                <span key={`joined-${member.id}`} className="card__eyebrow">
                  {formatTimestamp(member.joinedAt)}
                </span>,
                <div key={`actions-container-${member.id}`} className="dashboard-topbar__actions">
                  {hasManageAccess && member.userId !== user?.id && (
                    <Button
                      variant="ghost"
                      onClick={() => handleStatusChange(member, member.status === 'active' ? 'deactivated' : 'active')}
                    >
                      {member.status === 'active' ? 'Deactivate' : 'Activate'}
                    </Button>
                  )}
                  {hasRemoveAccess && member.userId !== user?.id && (
                    <Button
                      variant="ghost"
                      onClick={() => {
                        setTargetMember(member)
                        setConfirmRemoveOpen(true)
                      }}
                    >
                      Remove
                    </Button>
                  )}
                  {isOwner && member.userId !== user?.id && member.status === 'active' && (
                    <Button
                      variant="ghost"
                      onClick={() => {
                        setTargetMember(member)
                        setConfirmTransferOpen(true)
                      }}
                    >
                      Transfer Owner
                    </Button>
                  )}
                </div>
              ])}
            />

            {/* Pagination Controls */}
            <div className="page-actions">
              <span className="card__eyebrow">
                Showing Page {page} of {totalPages} ({totalCount} members)
              </span>
              <div className="dashboard-topbar__actions">
                <Button variant="secondary" onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}>
                  Previous
                </Button>
                <Button variant="secondary" onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}>
                  Next
                </Button>
              </div>
            </div>
          </div>
        )}
      </Card>

      {/* Invite Member Dialog */}
      <Dialog
        open={inviteOpen}
        title="Invite Team Member"
        description="Provide the email address and assign an initial access role for the new team member."
        onClose={() => setInviteOpen(false)}
      >
        <form onSubmit={handleInvite} className="stack">
          <Input
            label="Email Address"
            type="email"
            placeholder="colleague@company.com"
            value={inviteEmail}
            onChange={e => setInviteEmail(e.target.value)}
            required
            disabled={inviting}
          />

          <label className="field">
            <span className="field__label">Role</span>
            <select
              value={inviteRoleId}
              onChange={e => setInviteRoleId(e.target.value)}
              disabled={inviting}
              className="select"
            >
              {roles.map(r => (
                <option key={r.id} value={r.id}>
                  {r.name}
                </option>
              ))}
            </select>
          </label>

          <div className="page-actions">
            <Button variant="ghost" onClick={() => setInviteOpen(false)} disabled={inviting}>
              Cancel
            </Button>
            <Button variant="primary" type="submit" disabled={inviting || !inviteEmail.trim() || !inviteRoleId}>
              {inviting ? 'Inviting...' : 'Send Invitation'}
            </Button>
          </div>
        </form>
      </Dialog>

      {/* Remove Member Confirmation Dialog */}
      <Dialog
        open={confirmRemoveOpen}
        title="Remove Member Access"
        description={`Are you sure you want to remove ${targetMember?.username || targetMember?.email} from this organization? This action is immediate and will revoke all current access.`}
        onClose={() => setConfirmRemoveOpen(false)}
      >
        <div className="page-actions">
          <Button variant="ghost" onClick={() => setConfirmRemoveOpen(false)} disabled={removing}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleRemove} disabled={removing}>
            {removing ? 'Removing...' : 'Confirm Removal'}
          </Button>
        </div>
      </Dialog>

      {/* Transfer Ownership Confirmation Dialog */}
      <Dialog
        open={confirmTransferOpen}
        title="Transfer Organization Ownership"
        description={`Warning: This will transfer primary ownership of the organization to ${targetMember?.username || targetMember?.email}. You will lose ownership permissions and be demoted to their current role.`}
        onClose={() => setConfirmTransferOpen(false)}
      >
        <div className="page-actions">
          <Button variant="ghost" onClick={() => setConfirmTransferOpen(false)} disabled={transferring}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleTransferOwnership} disabled={transferring}>
            {transferring ? 'Transferring...' : 'Confirm Transfer'}
          </Button>
        </div>
      </Dialog>
    </div>
  )
}
