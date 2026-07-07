package rbac

import (
	"context"
	"errors"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

var (
	ErrAccessDenied     = errors.New("access denied")
	ErrResourceNotFound = errors.New("resource not found")
	ErrOrgNotFound      = errors.New("organization not found")
)

// Store defines persistent storage for RBAC data.
type Store interface {
	GetOrganization(ctx context.Context, orgID string) (*Organization, error)
	GetMembership(ctx context.Context, userID, orgID string) (*Membership, error)
	GetRole(ctx context.Context, roleID string) (*Role, error)
	GetResourcePolicy(ctx context.Context, resourceID, resourceType string) (*ResourcePolicy, error)
	UpsertOrganization(ctx context.Context, org *Organization) error
	UpsertRole(ctx context.Context, role *Role) error
	UpsertMembership(ctx context.Context, mem *Membership) error
	ProvisionOrganizationTx(ctx context.Context, org *Organization, role *Role, mem *Membership) error
	ListOrganizationsForUser(ctx context.Context, userID string, params common.QueryParams) ([]*Organization, int, error)
	ListMemberships(ctx context.Context, orgID string, params common.QueryParams) ([]*MembershipDetails, int, error)
	UpdateMembershipStatus(ctx context.Context, orgID, userID string, status string) error
	UpdateMembershipRole(ctx context.Context, orgID, userID string, roleID string) error
	RemoveMembership(ctx context.Context, orgID, userID string) error
	InviteMember(ctx context.Context, orgID string, email string, roleID string) error
	ListRoles(ctx context.Context, orgID string) ([]*Role, error)
}

// Service handles authorization evaluations.
type Service interface {
	Authorize(ctx context.Context, userID, orgID, permission string) error
	AuthorizeResource(ctx context.Context, userID, resourceID, resourceType, permission string) error
	CreateOrganization(ctx context.Context, org *Organization) error
	CreateRole(ctx context.Context, role *Role) error
	AddMembership(ctx context.Context, mem *Membership) error
	ProvisionOrganization(ctx context.Context, org *Organization, role *Role, mem *Membership) error
	ListOrganizations(ctx context.Context, userID string, params common.QueryParams) ([]*Organization, int, error)
	ListMembers(ctx context.Context, userID, orgID string, params common.QueryParams) ([]*MembershipDetails, int, error)
	InviteUser(ctx context.Context, adminID, orgID string, email string, roleID string) error
	UpdateMemberRole(ctx context.Context, adminID, orgID, targetUserID string, roleID string) error
	SetMemberStatus(ctx context.Context, adminID, orgID, targetUserID string, status string) error
	RemoveMember(ctx context.Context, adminID, orgID, targetUserID string) error
	TransferOrgOwnership(ctx context.Context, adminID, orgID, newOwnerID string) error
	ListRoles(ctx context.Context, orgID string) ([]*Role, error)
}

// DefaultService implements core RBAC authorization rules.
type DefaultService struct {
	log   zerolog.Logger
	store Store
}

func NewService(log zerolog.Logger, store Store) *DefaultService {
	return &DefaultService{
		log:   log.With().Str("component", "RBACService").Logger(),
		store: store,
	}
}

// Authorize verifies if a user holds a required permission in an organization.
func (s *DefaultService) Authorize(ctx context.Context, userID, orgID, permission string) error {
	membership, err := s.store.GetMembership(ctx, userID, orgID)
	if err != nil {
		s.log.Warn().Err(err).Str("user_id", userID).Str("org_id", orgID).Msg("Membership not found")
		return ErrAccessDenied
	}

	role, err := s.store.GetRole(ctx, membership.RoleID)
	if err != nil {
		s.log.Error().Err(err).Str("role_id", membership.RoleID).Msg("Role not found")
		return ErrAccessDenied
	}

	for _, p := range role.Permissions {
		if p == permission || p == "*" {
			return nil
		}
	}

	s.log.Debug().Str("user_id", userID).Str("org_id", orgID).Str("permission", permission).Msg("Access denied")
	return ErrAccessDenied
}

// AuthorizeResource verifies if a user holds a required permission for a specific resource via its owning organization.
func (s *DefaultService) AuthorizeResource(ctx context.Context, userID, resourceID, resourceType, permission string) error {
	policy, err := s.store.GetResourcePolicy(ctx, resourceID, resourceType)
	if err != nil {
		return ErrResourceNotFound
	}

	return s.Authorize(ctx, userID, policy.OrgID, permission)
}

func (s *DefaultService) CreateOrganization(ctx context.Context, org *Organization) error {
	return s.store.UpsertOrganization(ctx, org)
}

func (s *DefaultService) CreateRole(ctx context.Context, role *Role) error {
	return s.store.UpsertRole(ctx, role)
}

func (s *DefaultService) AddMembership(ctx context.Context, mem *Membership) error {
	return s.store.UpsertMembership(ctx, mem)
}

func (s *DefaultService) ProvisionOrganization(ctx context.Context, org *Organization, role *Role, mem *Membership) error {
	return s.store.ProvisionOrganizationTx(ctx, org, role, mem)
}

func (s *DefaultService) ListOrganizations(ctx context.Context, userID string, params common.QueryParams) ([]*Organization, int, error) {
	return s.store.ListOrganizationsForUser(ctx, userID, params)
}

func (s *DefaultService) ListMembers(ctx context.Context, userID, orgID string, params common.QueryParams) ([]*MembershipDetails, int, error) {
	if err := s.Authorize(ctx, userID, orgID, "members.view"); err != nil {
		return nil, 0, err
	}
	return s.store.ListMemberships(ctx, orgID, params)
}

func (s *DefaultService) InviteUser(ctx context.Context, adminID, orgID string, email string, roleID string) error {
	if err := s.Authorize(ctx, adminID, orgID, "members.invite"); err != nil {
		return err
	}
	return s.store.InviteMember(ctx, orgID, email, roleID)
}

func (s *DefaultService) UpdateMemberRole(ctx context.Context, adminID, orgID, targetUserID string, roleID string) error {
	if err := s.Authorize(ctx, adminID, orgID, "members.manage"); err != nil {
		return err
	}
	return s.store.UpdateMembershipRole(ctx, orgID, targetUserID, roleID)
}

func (s *DefaultService) SetMemberStatus(ctx context.Context, adminID, orgID, targetUserID string, status string) error {
	if err := s.Authorize(ctx, adminID, orgID, "members.manage"); err != nil {
		return err
	}
	return s.store.UpdateMembershipStatus(ctx, orgID, targetUserID, status)
}

func (s *DefaultService) RemoveMember(ctx context.Context, adminID, orgID, targetUserID string) error {
	if err := s.Authorize(ctx, adminID, orgID, "members.remove"); err != nil {
		return err
	}
	return s.store.RemoveMembership(ctx, orgID, targetUserID)
}

func (s *DefaultService) TransferOrgOwnership(ctx context.Context, adminID, orgID, newOwnerID string) error {
	if err := s.Authorize(ctx, adminID, orgID, "org.owner"); err != nil {
		return err
	}
	oldMem, err := s.store.GetMembership(ctx, adminID, orgID)
	if err != nil {
		return err
	}
	newMem, err := s.store.GetMembership(ctx, newOwnerID, orgID)
	if err != nil {
		return err
	}
	oldRoleID := oldMem.RoleID
	newRoleID := newMem.RoleID

	err = s.store.UpdateMembershipRole(ctx, orgID, newOwnerID, oldRoleID)
	if err != nil {
		return err
	}
	return s.store.UpdateMembershipRole(ctx, orgID, adminID, newRoleID)
}

func (s *DefaultService) ListRoles(ctx context.Context, orgID string) ([]*Role, error) {
	return s.store.ListRoles(ctx, orgID)
}
