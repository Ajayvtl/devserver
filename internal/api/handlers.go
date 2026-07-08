package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/Ajayvtl/devserver/internal/providerconfig"
	"github.com/Ajayvtl/devserver/internal/rbac"
	"github.com/Ajayvtl/devserver/internal/settings"
)

func (router *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
		return
	}

	providerName, ok := req["provider"].(string)
	if !ok || providerName == "" {
		providerName = "local" // fallback to default local provider
	}

	tokens, err := router.authService.Login(r.Context(), providerName, req)
	if err != nil {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	WriteSuccess(w, r, http.StatusOK, tokens)
}

func (router *Router) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
		return
	}

	tokens, err := router.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	WriteSuccess(w, r, http.StatusOK, tokens)
}

func (router *Router) handleOrganizations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	if r.Method == http.MethodGet {
		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		orgs, total, err := router.rbacService.ListOrganizations(r.Context(), userID, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve organizations", nil)
			return
		}

		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}

		WriteSuccessPaginated(w, r, http.StatusOK, orgs, meta)
		return
	} else if r.Method == http.MethodPost {
		var req rbac.Organization
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if req.Name == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Organization name is required", nil)
			return
		}

		role := &rbac.Role{
			Name: "owner",
			Permissions: []string{
				rbac.PermissionOrgOwner,
				rbac.PermissionMembersView,
				rbac.PermissionMembersInvite,
				rbac.PermissionMembersManage,
				rbac.PermissionMembersRemove,
				rbac.PermissionSettingsView,
				rbac.PermissionSettingsWrite,
				rbac.PermissionEnvironmentsView,
				rbac.PermissionEnvironmentsWrite,
				rbac.PermissionSecretsView,
				rbac.PermissionSecretsWrite,
				rbac.PermissionVariablesView,
				rbac.PermissionVariablesWrite,
				rbac.PermissionProvidersView,
				rbac.PermissionProvidersWrite,
				rbac.PermissionRolesView,
				rbac.PermissionOrganizationsView,
				rbac.PermissionOrganizationsCreate,
			},
		}
		mem := &rbac.Membership{
			UserID: userID,
		}

		if err := router.rbacService.ProvisionOrganization(r.Context(), &req, role, mem); err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to provision organization", nil)
			return
		}

		WriteSuccess(w, r, http.StatusCreated, req)
		return
	}

	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleUsersMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	orgID := r.Header.Get("X-Org-ID")
	roleName := "viewer" // default fallback
	var perms []string
	if orgID != "" {
		if role, err := router.rbacService.GetUserRole(r.Context(), userID, orgID); err == nil {
			roleName = role
		}
		if p, err := router.rbacService.GetUserPermissions(r.Context(), userID, orgID); err == nil {
			perms = p
		}
	} else {
		orgs, _, err := router.rbacService.ListOrganizations(r.Context(), userID, common.QueryParams{Page: 1, PerPage: 1})
		if err == nil && len(orgs) > 0 {
			if role, err := router.rbacService.GetUserRole(r.Context(), userID, orgs[0].ID); err == nil {
				roleName = role
			}
			if p, err := router.rbacService.GetUserPermissions(r.Context(), userID, orgs[0].ID); err == nil {
				perms = p
			}
		}
	}

	if perms == nil {
		perms = []string{}
	}

	WriteSuccess(w, r, http.StatusOK, map[string]any{
		"id":          userID,
		"status":      "active",
		"role":        roleName,
		"permissions": perms,
	})
}

func (router *Router) handleRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		orgID := r.Header.Get("X-Org-ID")
		if orgID == "" {
			WriteSuccess(w, r, http.StatusOK, []map[string]string{{"id": "role-1", "name": "admin"}})
			return
		}
		roles, err := router.rbacService.ListRoles(r.Context(), orgID)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve roles", nil)
			return
		}
		WriteSuccess(w, r, http.StatusOK, roles)
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scope := r.URL.Query().Get("scope")
		ownerID := r.URL.Query().Get("ownerId")
		if scope == "" || ownerID == "" {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "scope and ownerId are required", nil)
			return
		}

		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		sets, total, err := router.settingsService.ListSettings(r.Context(), settings.Scope(scope), ownerID, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve settings", nil)
			return
		}
		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}
		WriteSuccessPaginated(w, r, http.StatusOK, sets, meta)
		return
	} else if r.Method == http.MethodPost {
		var req settings.Setting
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.settingsService.SaveSetting(r.Context(), req.Scope, req.OwnerID, req.Key, req.Value); err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save setting", nil)
			return
		}
		WriteSuccess(w, r, http.StatusOK, req)
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleEnvironments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ownerID := r.URL.Query().Get("ownerId")
		if ownerID == "" {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "ownerId is required", nil)
			return
		}
		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		envs, total, err := router.envService.ListEnvironments(r.Context(), ownerID, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve environments", nil)
			return
		}

		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}

		WriteSuccessPaginated(w, r, http.StatusOK, envs, meta)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			OwnerID string `json:"ownerId"`
			Name    string `json:"name"`
			Type    string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if req.OwnerID == "" || req.Name == "" || req.Type == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Missing required fields", nil)
			return
		}
		env, err := router.envService.CreateEnvironment(r.Context(), req.OwnerID, req.Name, environments.EnvType(req.Type))
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create environment", nil)
			return
		}
		WriteSuccess(w, r, http.StatusCreated, env)
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			EnvID     string `json:"envId"`
			Key       string `json:"key"`
			Plaintext string `json:"plaintext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if err := router.envService.SetSecret(r.Context(), req.EnvID, req.Key, req.Plaintext); err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save secret", nil)
			return
		}
		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "saved"})
		return
	} else if r.Method == http.MethodGet {
		envID := r.URL.Query().Get("envId")
		if envID == "" {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "envId is required", nil)
			return
		}
		// We only want to return secret references, not raw values.
		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		secrets, total, err := router.envService.ListSecrets(r.Context(), envID, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve secrets", nil)
			return
		}
		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}
		WriteSuccessPaginated(w, r, http.StatusOK, secrets, meta)
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleVariables(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		envID := r.URL.Query().Get("envId")
		if envID == "" {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "envId is required", nil)
			return
		}
		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		vars, total, err := router.envService.ListVariables(r.Context(), envID, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve variables", nil)
			return
		}
		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}
		WriteSuccessPaginated(w, r, http.StatusOK, vars, meta)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			EnvID string `json:"envId"`
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.envService.SetVariable(r.Context(), req.EnvID, req.Key, req.Value); err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save variable", nil)
			return
		}
		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "saved"})
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scope := r.URL.Query().Get("scope")
		owner := r.URL.Query().Get("ownerId")
		params, err := common.ParseQueryParams(r, []string{"created_at", "name", "key"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}
		cfgs, total, err := router.providerService.ListConfigs(r.Context(), scope, owner, params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve provider configs", nil)
			return
		}
		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}
		WriteSuccessPaginated(w, r, http.StatusOK, cfgs, meta)
		return
	} else if r.Method == http.MethodPost {
		var req providerconfig.ProviderConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.providerService.SaveConfig(r.Context(), &req); err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save provider config", nil)
			return
		}
		WriteSuccess(w, r, http.StatusCreated, req)
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleProviderTest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			OwnerID   string `json:"ownerId"`
			Name      string `json:"name"`
			SecretRef string `json:"secretRef"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		ok, err := router.providerService.TestConnection(r.Context(), req.OwnerID, req.Name, req.SecretRef)
		if err != nil || !ok {
			msg := "Connection failed"
			if err != nil {
				msg = err.Error()
			}
			WriteError(w, r, http.StatusBadRequest, "TEST_FAILED", msg, nil)
			return
		}

		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	orgID := r.Header.Get("X-Org-ID")
	if orgID == "" {
		WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "X-Org-ID header required", nil)
		return
	}

	if r.Method == http.MethodGet {
		params, err := common.ParseQueryParams(r, []string{"username", "email", "status", "created_at"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}

		members, total, err := router.rbacService.ListMembers(r.Context(), userID, orgID, params)
		if err != nil {
			if errors.Is(err, rbac.ErrAccessDenied) {
				WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			} else {
				WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			}
			return
		}

		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}

		WriteSuccessPaginated(w, r, http.StatusOK, members, meta)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			Email  string `json:"email"`
			RoleID string `json:"roleId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if req.Email == "" || req.RoleID == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "email and roleId are required", nil)
			return
		}

		if err := router.rbacService.InviteUser(r.Context(), userID, orgID, req.Email, req.RoleID); err != nil {
			if errors.Is(err, rbac.ErrAccessDenied) {
				WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			} else {
				WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			}
			return
		}

		WriteSuccess(w, r, http.StatusCreated, map[string]string{"status": "invited"})
		return
	} else if r.Method == http.MethodPut {
		var req struct {
			UserID string `json:"userId"`
			RoleID string `json:"roleId"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if req.UserID == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "userId is required", nil)
			return
		}

		if req.RoleID != "" {
			if err := router.rbacService.UpdateMemberRole(r.Context(), userID, orgID, req.UserID, req.RoleID); err != nil {
				if errors.Is(err, rbac.ErrAccessDenied) {
					WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
				} else {
					WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
				}
				return
			}
		}

		if req.Status != "" {
			if err := router.rbacService.SetMemberStatus(r.Context(), userID, orgID, req.UserID, req.Status); err != nil {
				if errors.Is(err, rbac.ErrAccessDenied) {
					WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
				} else {
					WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
				}
				return
			}
		}

		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "updated"})
		return
	} else if r.Method == http.MethodDelete {
		targetUserID := r.URL.Query().Get("userId")
		if targetUserID == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "userId query parameter is required", nil)
			return
		}

		if err := router.rbacService.RemoveMember(r.Context(), userID, orgID, targetUserID); err != nil {
			if errors.Is(err, rbac.ErrAccessDenied) {
				WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			} else {
				WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			}
			return
		}

		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "removed"})
		return
	}

	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleTransferOwnership(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		WriteError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	orgID := r.Header.Get("X-Org-ID")
	if orgID == "" {
		WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "X-Org-ID header required", nil)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			NewOwnerUserID string `json:"newOwnerUserId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if req.NewOwnerUserID == "" {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "newOwnerUserId is required", nil)
			return
		}

		if err := router.rbacService.TransferOrgOwnership(r.Context(), userID, orgID, req.NewOwnerUserID); err != nil {
			if errors.Is(err, rbac.ErrAccessDenied) {
				WriteError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
			} else {
				WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			}
			return
		}

		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "ownership_transferred"})
		return
	}

	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleSuperAdminOrganizations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		params, err := common.ParseQueryParams(r, []string{"name", "created_at"})
		if err != nil {
			WriteError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
			return
		}

		orgs, total, err := router.rbacService.ListAllOrganizations(r.Context(), params)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error(), nil)
			return
		}

		totalPages := (total + params.PerPage - 1) / params.PerPage
		meta := PaginationMeta{
			Total:      total,
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: totalPages,
		}

		WriteSuccessPaginated(w, r, http.StatusOK, orgs, meta)
		return
	}

	WriteError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}
