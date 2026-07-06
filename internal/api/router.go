package api

import (
	"net/http"

	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/providerconfig"
	"github.com/Ajayvtl/devserver/internal/rbac"
	"github.com/Ajayvtl/devserver/internal/settings"
)

type Router struct {
	authService     auth.Service
	rbacService     rbac.Service
	settingsService settings.Service
	envService      environments.Service
	providerService providerconfig.Service
	bus             events.Bus
}

func NewRouter(
	authService auth.Service,
	rbacService rbac.Service,
	settingsService settings.Service,
	envService environments.Service,
	providerService providerconfig.Service,
	bus events.Bus,
) *Router {
	return &Router{
		authService:     authService,
		rbacService:     rbacService,
		settingsService: settingsService,
		envService:      envService,
		providerService: providerService,
		bus:             bus,
	}
}

func (router *Router) Register(mux *http.ServeMux) {
	// Group v1 API
	// ----------------------------------------------------
	// 1. Authentication
	mux.HandleFunc("/api/v1/auth/login", router.handleLogin)
	mux.HandleFunc("/api/v1/auth/refresh", router.handleRefresh)

	// Protected endpoints wrapper
	protect := func(h http.Handler) http.Handler {
		return IdempotencyMiddleware()(AuthMiddleware(router.authService)(AuditMiddleware(router.bus)(h)))
	}

	// 2. Organizations
	mux.Handle("/api/v1/organizations", protect(http.HandlerFunc(router.handleOrganizations)))

	// 3. Users
	mux.Handle("/api/v1/users/me", protect(http.HandlerFunc(router.handleUsersMe)))

	// 4. Roles
	mux.Handle("/api/v1/roles", protect(http.HandlerFunc(router.handleRoles)))

	// 5. Settings
	mux.Handle("/api/v1/settings", protect(http.HandlerFunc(router.handleSettings)))

	// 6. Environments
	mux.Handle("/api/v1/environments", protect(http.HandlerFunc(router.handleEnvironments)))

	// 7. Secrets and Variables
	mux.Handle("/api/v1/secrets", protect(http.HandlerFunc(router.handleSecrets)))
	mux.Handle("/api/v1/variables", protect(http.HandlerFunc(router.handleVariables)))

	// 8. Provider Configuration
	mux.Handle("/api/v1/providers", protect(http.HandlerFunc(router.handleProviders)))
	mux.Handle("/api/v1/providers/test", protect(http.HandlerFunc(router.handleProviderTest)))
}
