package rbac

import (
	"context"
	"testing"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

type mockStore struct {
	roles       map[string]*Role
	memberships map[string]*Membership
	calls       []string
}

func newMockStore() *mockStore {
	return &mockStore{
		roles:       map[string]*Role{},
		memberships: map[string]*Membership{},
	}
}

func (m *mockStore) GetOrganization(ctx context.Context, orgID string) (*Organization, error) {
	return &Organization{ID: orgID, Name: "mock"}, nil
}
func (m *mockStore) GetMembership(ctx context.Context, userID, orgID string) (*Membership, error) {
	for _, mem := range m.memberships {
		if mem.UserID == userID && mem.OrgID == orgID {
			return mem, nil
		}
	}
	return nil, ErrAccessDenied
}
func (m *mockStore) GetRole(ctx context.Context, roleID string) (*Role, error) {
	if r, ok := m.roles[roleID]; ok {
		return r, nil
	}
	return nil, ErrResourceNotFound
}
func (m *mockStore) GetResourcePolicy(ctx context.Context, resourceID, resourceType string) (*ResourcePolicy, error) {
	return nil, ErrResourceNotFound
}
func (m *mockStore) UpsertOrganization(ctx context.Context, org *Organization) error { return nil }
func (m *mockStore) UpsertRole(ctx context.Context, role *Role) error {
	m.roles[role.ID] = role
	return nil
}
func (m *mockStore) UpsertMembership(ctx context.Context, mem *Membership) error {
	m.memberships[mem.ID] = mem
	return nil
}
func (m *mockStore) ProvisionOrganizationTx(ctx context.Context, org *Organization, role *Role, mem *Membership) error {
	return nil
}
func (m *mockStore) ListOrganizationsForUser(ctx context.Context, userID string, params common.QueryParams) ([]*Organization, int, error) {
	return nil, 0, nil
}
func (m *mockStore) ListMemberships(ctx context.Context, orgID string, params common.QueryParams) ([]*MembershipDetails, int, error) {
	return nil, 0, nil
}
func (m *mockStore) UpdateMembershipStatus(ctx context.Context, orgID, userID string, status string) error {
	m.calls = append(m.calls, "status")
	return nil
}
func (m *mockStore) UpdateMembershipRole(ctx context.Context, orgID, userID string, roleID string) error {
	m.calls = append(m.calls, "role")
	return nil
}
func (m *mockStore) RemoveMembership(ctx context.Context, orgID, userID string) error {
	m.calls = append(m.calls, "remove")
	return nil
}
func (m *mockStore) InviteMember(ctx context.Context, orgID string, email string, roleID string) error {
	m.calls = append(m.calls, "invite")
	return nil
}
func (m *mockStore) ListRoles(ctx context.Context, orgID string) ([]*Role, error) { return nil, nil }
func (m *mockStore) ListAllOrganizations(ctx context.Context, params common.QueryParams) ([]*Organization, int, error) {
	return nil, 0, nil
}

func TestAuthorizeWildcardAndDeny(t *testing.T) {
	ms := newMockStore()
	// role with wildcard
	r1 := &Role{ID: "r1", Name: "owner", Permissions: []string{PermissionWildcard}}
	ms.roles[r1.ID] = r1
	mem := &Membership{ID: "m1", OrgID: "org1", UserID: "user1", RoleID: r1.ID}
	ms.memberships[mem.ID] = mem

	svc := NewService(zerolog.Nop(), ms)

	if err := svc.Authorize(context.Background(), "user1", "org1", PermissionMembersView); err != nil {
		t.Fatalf("expected wildcard to authorize, got %v", err)
	}

	// role without permission
	r2 := &Role{ID: "r2", Name: "limited", Permissions: []string{"environments.view"}}
	ms.roles[r2.ID] = r2
	mem2 := &Membership{ID: "m2", OrgID: "org1", UserID: "user2", RoleID: r2.ID}
	ms.memberships[mem2.ID] = mem2

	if err := svc.Authorize(context.Background(), "user2", "org1", PermissionMembersView); err == nil {
		t.Fatalf("expected access denied for missing permission")
	}
}

func TestGetUserPermissionsAndTransfer(t *testing.T) {
	ms := newMockStore()
	rOwner := &Role{ID: "r-owner", Name: "owner", Permissions: []string{PermissionOrgOwner, PermissionMembersManage}}
	rViewer := &Role{ID: "r-viewer", Name: "viewer", Permissions: []string{PermissionMembersView}}
	ms.roles[rOwner.ID] = rOwner
	ms.roles[rViewer.ID] = rViewer

	// admin user is owner
	memA := &Membership{ID: "mA", OrgID: "orgX", UserID: "admin", RoleID: rOwner.ID}
	ms.memberships[memA.ID] = memA
	// bob is viewer
	memB := &Membership{ID: "mB", OrgID: "orgX", UserID: "bob", RoleID: rViewer.ID}
	ms.memberships[memB.ID] = memB

	svc := NewService(zerolog.Nop(), ms)

	// transfer ownership: admin must have PermissionOrgOwner
	if err := svc.TransferOrgOwnership(context.Background(), "admin", "orgX", "bob"); err != nil {
		t.Fatalf("expected transfer to succeed, got %v", err)
	}

	// ensure UpdateMembershipRole was called twice (for new owner and demotion)
	found := 0
	for _, c := range ms.calls {
		if c == "role" {
			found++
		}
	}
	if found < 2 {
		t.Fatalf("expected membership role updates, got calls=%v", ms.calls)
	}
}
