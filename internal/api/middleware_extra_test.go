package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/rbac"
)

type mockRBACService2 struct {
	lastUser string
	lastOrg  string
	lastPerm string
	allow    bool
}

func (m *mockRBACService2) Authorize(ctx context.Context, userID, orgID, permission string) error {
	m.lastUser = userID
	m.lastOrg = orgID
	m.lastPerm = permission
	if m.allow {
		return nil
	}
	return rbac.ErrAccessDenied
}
func (m *mockRBACService2) AuthorizeResource(ctx context.Context, userID, resourceID, resourceType, permission string) error {
	return rbac.ErrAccessDenied
}
func (m *mockRBACService2) CreateOrganization(ctx context.Context, org *rbac.Organization) error {
	return nil
}
func (m *mockRBACService2) CreateRole(ctx context.Context, role *rbac.Role) error         { return nil }
func (m *mockRBACService2) AddMembership(ctx context.Context, mem *rbac.Membership) error { return nil }
func (m *mockRBACService2) ProvisionOrganization(ctx context.Context, org *rbac.Organization, role *rbac.Role, mem *rbac.Membership) error {
	return nil
}
func (m *mockRBACService2) ListOrganizations(ctx context.Context, userID string, params common.QueryParams) ([]*rbac.Organization, int, error) {
	return nil, 0, nil
}
func (m *mockRBACService2) ListMembers(ctx context.Context, userID, orgID string, params common.QueryParams) ([]*rbac.MembershipDetails, int, error) {
	return nil, 0, nil
}
func (m *mockRBACService2) InviteUser(ctx context.Context, adminID, orgID string, email string, roleID string) error {
	return nil
}
func (m *mockRBACService2) UpdateMemberRole(ctx context.Context, adminID, orgID, targetUserID string, roleID string) error {
	return nil
}
func (m *mockRBACService2) SetMemberStatus(ctx context.Context, adminID, orgID, targetUserID string, status string) error {
	return nil
}
func (m *mockRBACService2) RemoveMember(ctx context.Context, adminID, orgID, targetUserID string) error {
	return nil
}
func (m *mockRBACService2) TransferOrgOwnership(ctx context.Context, adminID, orgID, newOwnerID string) error {
	return nil
}
func (m *mockRBACService2) ListRoles(ctx context.Context, orgID string) ([]*rbac.Role, error) {
	return nil, nil
}
func (m *mockRBACService2) GetUserRole(ctx context.Context, userID, orgID string) (string, error) {
	return "", nil
}
func (m *mockRBACService2) GetUserPermissions(ctx context.Context, userID, orgID string) ([]string, error) {
	return nil, nil
}
func (m *mockRBACService2) ListAllOrganizations(ctx context.Context, params common.QueryParams) ([]*rbac.Organization, int, error) {
	return nil, 0, nil
}

func TestMethodRBACMiddleware_PATCH_and_unmapped(t *testing.T) {
	svc := &mockRBACService2{allow: true}
	mw := MethodRBACMiddleware(svc, map[string]string{"PATCH": rbac.PermissionMembersManage, "POST": rbac.PermissionMembersInvite})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	// test PATCH mapping invokes authorize
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/organizations/members", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, "u1"))
	req.Header.Set("X-Org-ID", "org1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if svc.lastPerm != rbac.PermissionMembersManage {
		t.Fatalf("expected permission %s, got %s", rbac.PermissionMembersManage, svc.lastPerm)
	}

	// test HEAD (unmapped) passes through (no authorize call)
	svc.lastPerm = ""
	req2 := httptest.NewRequest(http.MethodHead, "/api/v1/organizations", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), ContextKeyUserID, "u1"))
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)
	if svc.lastPerm != "" {
		t.Fatalf("expected no authorize for unmapped method, got %s", svc.lastPerm)
	}
}

func TestAuditMiddleware_Publish(t *testing.T) {
	// simple bus capture
	published := 0
	bus := &testBus{onPublish: func(tpe string, payload any) {
		if tpe == "audit.mutation" {
			published++
		}
	}}

	mw := AuditMiddleware(bus)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) }))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/some", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, "u1"))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if published != 1 {
		t.Fatalf("expected 1 audit publish, got %d", published)
	}
}

// lightweight test bus
type testBus struct{ onPublish func(string, any) }

func (t *testBus) Publish(eventType events.EventType, payload any) {
	if t.onPublish != nil {
		t.onPublish(string(eventType), payload)
	}
}
func (t *testBus) Subscribe(eventType events.EventType) events.Subscriber {
	ch := make(chan events.Event, 1)
	return ch
}
