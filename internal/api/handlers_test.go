package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ajayvtl/devserver/internal/rbac"
)

type mockRBACForHandlers struct {
	rbac.Service
}

func (m *mockRBACForHandlers) ListOrganizations(ctx context.Context, userID string) ([]*rbac.Organization, error) {
	return []*rbac.Organization{{ID: "org-1", Name: "Test Org"}}, nil
}

func (m *mockRBACForHandlers) CreateOrganization(ctx context.Context, org *rbac.Organization) error {
	org.ID = "new-org-id"
	return nil
}

func (m *mockRBACForHandlers) CreateRole(ctx context.Context, role *rbac.Role) error {
	role.ID = "new-role-id"
	return nil
}

func (m *mockRBACForHandlers) AddMembership(ctx context.Context, mem *rbac.Membership) error {
	return nil
}

func TestHandleOrganizationsGet(t *testing.T) {
	router := NewRouter(nil, &mockRBACForHandlers{}, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/organizations", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, "user-1"))
	w := httptest.NewRecorder()

	router.handleOrganizations(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var res APIResponse
	json.NewDecoder(w.Body).Decode(&res)
	if !res.Success {
		t.Error("Expected success=true")
	}
}

func TestHandleOrganizationsPost(t *testing.T) {
	router := NewRouter(nil, &mockRBACForHandlers{}, nil, nil, nil)

	body, _ := json.Marshal(map[string]string{"name": "New Org"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, "user-1"))
	w := httptest.NewRecorder()

	router.handleOrganizations(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", w.Code)
	}

	var res APIResponse
	json.NewDecoder(w.Body).Decode(&res)
	if !res.Success {
		t.Error("Expected success=true")
	}
}
