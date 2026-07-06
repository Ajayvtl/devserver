package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ajayvtl/devserver/internal/rbac"
)

type mockTenantRBACService struct {
	rbac.Service
}

func (m *mockTenantRBACService) Authorize(ctx context.Context, userID, orgID, permission string) error {
	if userID == "user1" && orgID == "org2" {
		return rbac.ErrAccessDenied
	}
	return nil
}

func TestRBAC_DenyByDefault(t *testing.T) {
	// Directly test RBACMiddleware
	mockSvc := &mockTenantRBACService{}
	middleware := RBACMiddleware(mockSvc, "read")
	
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil)
	w := httptest.NewRecorder()
	
	// Call without auth context
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for missing auth context, got %d", w.Code)
	}
}

func TestRBAC_TenantIsolation(t *testing.T) {
	mockSvc := &mockTenantRBACService{}
	middleware := RBACMiddleware(mockSvc, "read")
	
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// User from org1 trying to access org2
	req := httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, "user1")
	req = req.WithContext(ctx)
	req.Header.Set("X-Org-ID", "org2")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden for cross-tenant access, got %d", w.Code)
	}
}
