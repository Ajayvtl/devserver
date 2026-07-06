package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/rbac"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyOrgID  contextKey = "org_id"
	ContextKeyReqID  contextKey = "request_id"
)

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			w.Header().Set("X-Request-ID", reqID)
			ctx := context.WithValue(r.Context(), ContextKeyReqID, reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Missing Authorization header", nil)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Authorization header format", nil)
				return
			}

			userID, err := authService.ValidateToken(r.Context(), parts[1])
			if err != nil {
				WriteError(w, r, http.StatusUnauthorized, "SESSION_EXPIRED", "Token invalid or expired", nil)
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
				WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
				return
			}
			orgID := r.Header.Get("X-Org-ID") // Typical way to pass current tenant
			if orgID == "" {
				WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "X-Org-ID header required", nil)
				return
			}

			err := rbacService.Authorize(r.Context(), userID, orgID, requiredPermission)
			if err != nil {
				WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for this organization", nil)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyOrgID, orgID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func AuditMiddleware(bus events.Bus) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only audit mutations
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete && r.Method != http.MethodPatch {
				next.ServeHTTP(w, r)
				return
			}

			rr := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rr, r)

			if rr.status >= 200 && rr.status < 400 {
				userID, _ := r.Context().Value(ContextKeyUserID).(string)
				if userID == "" {
					userID = "system"
				}

				action := r.Method + " " + r.URL.Path

				// In a full implementation, we'd save this to DB directly or via bus
				// For WP-8.5 Audit Trail requirement:
				bus.Publish(events.EventType("audit.mutation"), map[string]any{
					"userId":    userID,
					"action":    action,
					"ipAddress": r.RemoteAddr,
					"userAgent": r.UserAgent(),
				})
			}
		})
	}
}
