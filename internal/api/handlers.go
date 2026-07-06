package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/Ajayvtl/devserver/internal/providerconfig"
	"github.com/Ajayvtl/devserver/internal/rbac"
	"github.com/Ajayvtl/devserver/internal/settings"
)

func (router *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
		return
	}

	providerName, ok := req["provider"].(string)
	if !ok || providerName == "" {
		providerName = "local" // fallback to default local provider
	}

	tokens, err := router.authService.Login(r.Context(), providerName, req)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	WriteSuccess(w, http.StatusOK, tokens)
}

func (router *Router) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
		return
	}

	tokens, err := router.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	WriteSuccess(w, http.StatusOK, tokens)
}

func (router *Router) handleOrganizations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	if r.Method == http.MethodGet {
		orgs, err := router.rbacService.ListOrganizations(r.Context(), userID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve organizations", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, orgs)
		return
	} else if r.Method == http.MethodPost {
		var req rbac.Organization
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if req.Name == "" {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Organization name is required", nil)
			return
		}

		if err := router.rbacService.CreateOrganization(r.Context(), &req); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create organization", nil)
			return
		}

		// Also make the creator an admin
		role := &rbac.Role{
			OrgID:       req.ID,
			Name:        "admin",
			Permissions: []string{"*"},
		}
		if err := router.rbacService.CreateRole(r.Context(), role); err == nil {
			_ = router.rbacService.AddMembership(r.Context(), &rbac.Membership{
				UserID: userID,
				OrgID:  req.ID,
				RoleID: role.ID,
			})
		}

		WriteSuccess(w, http.StatusCreated, req)
		return
	}

	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleUsersMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
		return
	}

	WriteSuccess(w, http.StatusOK, map[string]string{"id": userID, "status": "active"})
}

func (router *Router) handleRoles(w http.ResponseWriter, r *http.Request) {
	// Simple stub mapping for WP-8.1 REST endpoints wrapper.
	if r.Method == http.MethodGet {
		WriteSuccess(w, http.StatusOK, []map[string]string{{"id": "role-1", "name": "admin"}})
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scope := r.URL.Query().Get("scope")
		ownerID := r.URL.Query().Get("ownerId")
		if scope == "" || ownerID == "" {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "scope and ownerId are required", nil)
			return
		}

		sets, err := router.settingsService.ListSettings(r.Context(), settings.Scope(scope), ownerID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve settings", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, sets)
		return
	} else if r.Method == http.MethodPost {
		var req settings.Setting
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.settingsService.SaveSetting(r.Context(), req.Scope, req.OwnerID, req.Key, req.Value); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save setting", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, req)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleEnvironments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ownerID := r.URL.Query().Get("ownerId")
		if ownerID == "" {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "ownerId is required", nil)
			return
		}
		envs, err := router.envService.ListEnvironments(r.Context(), ownerID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve environments", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, envs)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			OwnerID string `json:"ownerId"`
			Name    string `json:"name"`
			Type    string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if req.OwnerID == "" || req.Name == "" || req.Type == "" {
			WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Missing required fields", nil)
			return
		}
		env, err := router.envService.CreateEnvironment(r.Context(), req.OwnerID, req.Name, environments.EnvType(req.Type))
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create environment", nil)
			return
		}
		WriteSuccess(w, http.StatusCreated, env)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			EnvID     string `json:"envId"`
			Key       string `json:"key"`
			Plaintext string `json:"plaintext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}

		if err := router.envService.SetSecret(r.Context(), req.EnvID, req.Key, req.Plaintext); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save secret", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, map[string]string{"status": "saved"})
		return
	} else if r.Method == http.MethodGet {
		envID := r.URL.Query().Get("envId")
		if envID == "" {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "envId is required", nil)
			return
		}
		// We only want to return secret references, not raw values.
		secrets, err := router.envService.ListSecrets(r.Context(), envID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve secrets", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, secrets)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleVariables(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		envID := r.URL.Query().Get("envId")
		if envID == "" {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "envId is required", nil)
			return
		}
		vars, err := router.envService.ListVariables(r.Context(), envID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve variables", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, vars)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			EnvID string `json:"envId"`
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.envService.SetVariable(r.Context(), req.EnvID, req.Key, req.Value); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save variable", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, map[string]string{"status": "saved"})
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func (router *Router) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scope := r.URL.Query().Get("scope")
		owner := r.URL.Query().Get("ownerId")
		cfgs, err := router.providerService.ListConfigs(r.Context(), scope, owner)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve provider configs", nil)
			return
		}
		WriteSuccess(w, http.StatusOK, cfgs)
		return
	} else if r.Method == http.MethodPost {
		var req providerconfig.ProviderConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload", nil)
			return
		}
		if err := router.providerService.SaveConfig(r.Context(), &req); err != nil {
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save provider config", nil)
			return
		}
		WriteSuccess(w, http.StatusCreated, req)
		return
	}
	WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}
