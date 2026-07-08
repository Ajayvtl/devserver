package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/rbac"
)

type mockRBACAuthService struct {
	auth.Service
}

func (m *mockRBACAuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	return token, nil
}

type mockTenantRBACService struct {
	rbac.Service
	orgs []*rbac.Organization
}

func (m *mockTenantRBACService) Authorize(ctx context.Context, userID, orgID, permission string) error {
	// Reject user1 accessing org2 resources (Tenant Isolation check)
	if userID == "user1" && orgID == "org2" {
		return rbac.ErrAccessDenied
	}

	// SuperAdmin permission check
	if permission == rbac.PermissionOrgAdmin {
		if userID != "superadmin-user" {
			return rbac.ErrAccessDenied
		}
	}

	// Membership check
	if userID == "unauthorized-user" {
		return rbac.ErrAccessDenied
	}

	return nil
}

func (m *mockTenantRBACService) GetUserRole(ctx context.Context, userID, orgID string) (string, error) {
	if userID == "superadmin-user" {
		return "super admin", nil
	}
	if userID == "user1" {
		return "admin", nil
	}
	return "viewer", nil
}

func (m *mockTenantRBACService) ListAllOrganizations(ctx context.Context, params common.QueryParams) ([]*rbac.Organization, int, error) {
	var filtered []*rbac.Organization
	for _, o := range m.orgs {
		if params.Query != "" {
			if !strings.Contains(strings.ToLower(o.Name), strings.ToLower(params.Query)) {
				continue
			}
		}
		filtered = append(filtered, o)
	}

	// Mock sorting
	if params.Sort == "name" {
		// Bubble sort name
		for i := 0; i < len(filtered); i++ {
			for j := i + 1; j < len(filtered); j++ {
				shouldSwap := false
				if params.Dir == "ASC" && filtered[i].Name > filtered[j].Name {
					shouldSwap = true
				} else if params.Dir == "DESC" && filtered[i].Name < filtered[j].Name {
					shouldSwap = true
				}
				if shouldSwap {
					filtered[i], filtered[j] = filtered[j], filtered[i]
				}
			}
		}
	}

	total := len(filtered)
	start := (params.Page - 1) * params.PerPage
	if start < 0 {
		start = 0
	}
	if start > len(filtered) {
		return []*rbac.Organization{}, total, nil
	}
	end := start + params.PerPage
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func TestRBAC_DenyByDefault(t *testing.T) {
	mockSvc := &mockTenantRBACService{}
	middleware := RBACMiddleware(mockSvc, "read")

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/environments", nil)
	w := httptest.NewRecorder()

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

	// User1 trying to access org2 (which is denied in mock tenant Service)
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

func TestSuperAdmin_OrganizationsEndpoint(t *testing.T) {
	// Create mock dataset (large dataset verification)
	var mockOrgs []*rbac.Organization
	for i := 1; i <= 50; i++ {
		mockOrgs = append(mockOrgs, &rbac.Organization{
			ID:        fmt.Sprintf("org-%d", i),
			Name:      fmt.Sprintf("Tenant Organization %03d", i),
			Slug:      fmt.Sprintf("tenant-%d", i),
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Hour),
		})
	}

	mockSvc := &mockTenantRBACService{orgs: mockOrgs}
	mockAuth := &mockRBACAuthService{}
	router := NewRouter(mockAuth, mockSvc, nil, nil, nil, nil)
	mux := http.NewServeMux()
	router.Register(mux)

	// Test 1: Access Denied for Unauthorized Role
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/superadmin/organizations", nil)
		req.Header.Set("Authorization", "Bearer unauthorized-user")
		req.Header.Set("X-Org-ID", "org-1")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden for unauthorized superadmin role, got %d", w.Code)
		}
	}

	// Test 2: Successful access with Pagination
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/superadmin/organizations?page=2&per_page=10", nil)
		req.Header.Set("Authorization", "Bearer superadmin-user")
		req.Header.Set("X-Org-ID", "org-1")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", w.Code)
		}

		var res struct {
			Success bool                 `json:"success"`
			Data    []*rbac.Organization `json:"data"`
			Meta    PaginationMeta       `json:"meta"`
		}
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !res.Success {
			t.Error("Expected success to be true")
		}
		if len(res.Data) != 10 {
			t.Errorf("Expected 10 items on page 2, got %d", len(res.Data))
		}
		if res.Meta.Total != 50 {
			t.Errorf("Expected total count of 50, got %d", res.Meta.Total)
		}
		if res.Meta.Page != 2 {
			t.Errorf("Expected current page meta 2, got %d", res.Meta.Page)
		}
	}

	// Test 3: Filtering & Searching
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/superadmin/organizations?q=Tenant%20Organization%20005", nil)
		req.Header.Set("Authorization", "Bearer superadmin-user")
		req.Header.Set("X-Org-ID", "org-1")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var res struct {
			Success bool                 `json:"success"`
			Data    []*rbac.Organization `json:"data"`
		}
		json.NewDecoder(w.Body).Decode(&res)

		if len(res.Data) != 1 {
			t.Errorf("Expected 1 filtered match, got %d", len(res.Data))
		}
		if res.Data[0].ID != "org-5" {
			t.Errorf("Expected filtered org to be org-5, got %s", res.Data[0].ID)
		}
	}

	// Test 4: Sorting (Name descending)
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/superadmin/organizations?sort=name&dir=DESC&per_page=5", nil)
		req.Header.Set("Authorization", "Bearer superadmin-user")
		req.Header.Set("X-Org-ID", "org-1")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		var res struct {
			Success bool                 `json:"success"`
			Data    []*rbac.Organization `json:"data"`
		}
		json.NewDecoder(w.Body).Decode(&res)

		if len(res.Data) != 5 {
			t.Fatalf("Expected 5 sorted items, got %d", len(res.Data))
		}
		// First item should be Tenant Organization 050 when sorted DESC
		if res.Data[0].Name != "Tenant Organization 050" {
			t.Errorf("Expected first element to be Tenant Organization 050 under DESC sort, got %s", res.Data[0].Name)
		}
	}
}

// helper func for naming formatting in formatting strings
func fmtOrgName(i int) string {
	return "Tenant Organization " + string(rune(i))
}
