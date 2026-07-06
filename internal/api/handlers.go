package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/Ajayvtl/devserver/internal/providerconfig"
	"github.com/Ajayvtl/devserver/internal/settings"
)

func (router *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	providerName, ok := req["provider"].(string)
	if !ok || providerName == "" {
		providerName = "local" // fallback to default local provider
	}

	tokens, err := router.authService.Login(r.Context(), providerName, req)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err)
		return
	}

	writeJSON(w, tokens)
}

func (router *Router) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := router.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err)
		return
	}

	writeJSON(w, tokens)
}

func (router *Router) handleOrganizations(w http.ResponseWriter, r *http.Request) {
	// Simple stub mapping for WP-8.1 REST endpoints wrapper.
	// We'd delegate to a CreateOrganization / ListOrganizations method in rbacService if it were public.
	// For now, we simulate the HTTP response structure.
	if r.Method == http.MethodGet {
		writeJSON(w, []map[string]string{{"id": "org-1", "name": "Default Org"}})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (router *Router) handleUsersMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, nil)
		return
	}

	writeJSON(w, map[string]string{"id": userID, "status": "active"})
}

func (router *Router) handleRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, []map[string]string{{"id": "role-1", "name": "admin"}})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (router *Router) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Example of wiring
		val, err := router.settingsService.GetSettingValue(r.Context(), "global", "system", "theme")
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, map[string]string{"theme": val})
		return
	} else if r.Method == http.MethodPost {
		var req settings.Setting
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		if err := router.settingsService.SaveSetting(r.Context(), req.Scope, req.OwnerID, req.Key, req.Value); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, req)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (router *Router) handleEnvironments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Typical query parameter usage
		ownerID := r.URL.Query().Get("ownerId")
		_ = ownerID
		envs, _ := router.envService.(*environments.DefaultService) // If we needed to cast for specific queries not in interface
		_ = envs
		// Returning a simulated response if list isn't directly exposed in the main envService interface
		// (The main envService interface in WP-7.4 doesn't have ListEnvironments, only Store does.
		// We'd expose it properly in Service in a real scenario.)
		writeJSON(w, []map[string]string{})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (router *Router) handleSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			EnvID     string `json:"envId"`
			Key       string `json:"key"`
			Plaintext string `json:"plaintext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}

		if err := router.envService.SetSecret(r.Context(), req.EnvID, req.Key, req.Plaintext); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, map[string]string{"status": "saved"})
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (router *Router) handleProviders(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		scope := r.URL.Query().Get("scope")
		owner := r.URL.Query().Get("ownerId")
		cfgs, err := router.providerService.ListConfigs(r.Context(), scope, owner)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, cfgs)
		return
	} else if r.Method == http.MethodPost {
		var req providerconfig.ProviderConfig
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		if err := router.providerService.SaveConfig(r.Context(), &req); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, req)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
