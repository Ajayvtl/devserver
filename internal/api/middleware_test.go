package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/rbac"
)

type mockAuthService struct {
	auth.Service
	userID string
	err    error
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.userID, nil
}

type mockRBACService struct {
	rbac.Service
	err error
}

func (m *mockRBACService) Authorize(ctx context.Context, userID, orgID, permission string) error {
	return m.err
}

func TestAuthMiddleware(t *testing.T) {
	mockAuth := &mockAuthService{userID: "user-1"}
	handler := AuthMiddleware(mockAuth)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing header, got %d", w.Code)
	}

	req.Header.Set("Authorization", "Bearer valid-token")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for valid token, got %d", w.Code)
	}
}

func TestRBACMiddleware(t *testing.T) {
	mockRBAC := &mockRBACService{}
	handler := RBACMiddleware(mockRBAC, "write")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), ContextKeyUserID, "user-1")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing X-Org-ID, got %d", w.Code)
	}

	req.Header.Set("X-Org-ID", "org-1")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 for valid org and permission, got %d", w.Code)
	}
}
