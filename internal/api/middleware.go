package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/rbac"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyOrgID  contextKey = "org_id"
)

func AuthMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing Authorization header", nil)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Authorization header format", nil)
				return
			}

			userID, err := authService.ValidateToken(r.Context(), parts[1])
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "SESSION_EXPIRED", "Token invalid or expired", nil)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RBACMiddleware(rbacService rbac.Service, requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(ContextKeyUserID).(string)
			if !ok || userID == "" {
				WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
				return
			}
			orgID := r.Header.Get("X-Org-ID") // Typical way to pass current tenant
			if orgID == "" {
				WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "X-Org-ID header required", nil)
				return
			}

			err := rbacService.Authorize(r.Context(), userID, orgID, requiredPermission)
			if err != nil {
				WriteError(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for this organization", nil)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyOrgID, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
