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

	// RBAC protected wrappers
	protectWithRBAC := func(requiredPermission string, h http.Handler) http.Handler {
		return IdempotencyMiddleware()(AuthMiddleware(router.authService)(RBACMiddleware(router.rbacService, requiredPermission)(AuditMiddleware(router.bus)(h))))
	}

	protectWithMethodRBAC := func(methodPermissions map[string]string, h http.Handler) http.Handler {
		return IdempotencyMiddleware()(AuthMiddleware(router.authService)(MethodRBACMiddleware(router.rbacService, methodPermissions)(AuditMiddleware(router.bus)(h))))
	}

	// 2. Organizations
	mux.Handle("/api/v1/organizations", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionOrganizationsView,
		"POST": rbac.PermissionOrganizationsCreate,
	}, http.HandlerFunc(router.handleOrganizations)))

	mux.Handle("/api/v1/organizations/members", protectWithMethodRBAC(map[string]string{
		"GET":    rbac.PermissionMembersView,
		"POST":   rbac.PermissionMembersInvite,
		"PUT":    rbac.PermissionMembersManage,
		"PATCH":  rbac.PermissionMembersManage,
		"DELETE": rbac.PermissionMembersRemove,
	}, http.HandlerFunc(router.handleOrganizationMembers)))

	mux.Handle("/api/v1/organizations/transfer-ownership", protectWithRBAC(rbac.PermissionOrgOwner, http.HandlerFunc(router.handleTransferOwnership)))
	// 3. Users
	mux.Handle("/api/v1/users/me", protect(http.HandlerFunc(router.handleUsersMe)))

	// 4. Roles
	mux.Handle("/api/v1/roles", protect(http.HandlerFunc(router.handleRoles)))

	// 5. Settings
	mux.Handle("/api/v1/settings", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionSettingsView,
		"POST": rbac.PermissionSettingsWrite,
	}, http.HandlerFunc(router.handleSettings)))

	// 6. Environments
	mux.Handle("/api/v1/environments", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionEnvironmentsView,
		"POST": rbac.PermissionEnvironmentsWrite,
	}, http.HandlerFunc(router.handleEnvironments)))

	// 7. Secrets and Variables
	mux.Handle("/api/v1/secrets", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionSecretsView,
		"POST": rbac.PermissionSecretsWrite,
	}, http.HandlerFunc(router.handleSecrets)))

	mux.Handle("/api/v1/variables", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionVariablesView,
		"POST": rbac.PermissionVariablesWrite,
	}, http.HandlerFunc(router.handleVariables)))

	// 8. Provider Configuration
	mux.Handle("/api/v1/providers", protectWithMethodRBAC(map[string]string{
		"GET":  rbac.PermissionProvidersView,
		"POST": rbac.PermissionProvidersWrite,
	}, http.HandlerFunc(router.handleProviders)))

	mux.Handle("/api/v1/providers/test", protectWithRBAC(rbac.PermissionProvidersWrite, http.HandlerFunc(router.handleProviderTest)))

	// 9. Super Admin Dashboard Endpoint
	mux.Handle("/api/v1/superadmin/organizations", protectWithRBAC(rbac.PermissionOrgAdmin, http.HandlerFunc(router.handleSuperAdminOrganizations)))
}
