package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Ajayvtl/devserver/internal/rbac"
)

type threadSafeMockRBAC struct {
	rbac.Service
	mu   sync.Mutex
	orgs map[string]*rbac.Organization
}

func (m *threadSafeMockRBAC) ProvisionOrganization(ctx context.Context, org *rbac.Organization, role *rbac.Role, mem *rbac.Membership) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	org.ID = "new-org-id"
	m.orgs[org.ID] = org
	return nil
}

func TestHandleOrganizationsConcurrency(t *testing.T) {
	mockRbac := &threadSafeMockRBAC{orgs: make(map[string]*rbac.Organization)}
	router := NewRouter(nil, mockRbac, nil, nil, nil, nil)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]string{"name": "Parallel Org"})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewBuffer(body))
			req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, "user-1"))
			w := httptest.NewRecorder()
			router.handleOrganizations(w, req)
			if w.Code != http.StatusCreated {
				t.Errorf("Expected 201 Created, got %d", w.Code)
			}
		}(i)
	}
	wg.Wait()
}
