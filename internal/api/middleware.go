package api

import (
	"context"
	"encoding/json"
	"errors"
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

func writeJSONError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func AuthMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, auth.ErrSessionExpired) // reusing error text for unauthorized
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeJSONError(w, http.StatusUnauthorized, auth.ErrSessionExpired)
				return
			}

			userID, err := authService.ValidateToken(r.Context(), parts[1])
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, auth.ErrSessionExpired)
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
				writeJSONError(w, http.StatusUnauthorized, auth.ErrSessionExpired)
				return
			}
			orgID := r.Header.Get("X-Org-ID") // Typical way to pass current tenant
			if orgID == "" {
				writeJSONError(w, http.StatusBadRequest, errors.New("X-Org-ID header required"))
				return
			}

			err := rbacService.Authorize(r.Context(), userID, orgID, requiredPermission)
			if err != nil {
				writeJSONError(w, http.StatusForbidden, errors.New("insufficient permissions"))
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyOrgID, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
