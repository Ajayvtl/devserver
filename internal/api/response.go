package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Success   bool   `json:"success"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   any    `json:"details,omitempty"`
	Data      any    `json:"data,omitempty"`
	Meta      any    `json:"meta,omitempty"`
	RequestId string `json:"requestId,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalPages int `json:"totalPages"`
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string, details any) {
	reqID := r.Header.Get("X-Request-ID")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success:   false,
		Code:      code,
		Message:   message,
		Details:   details,
		RequestId: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func WriteSuccess(w http.ResponseWriter, r *http.Request, status int, data any) {
	reqID := r.Header.Get("X-Request-ID")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success:   true,
		Data:      data,
		RequestId: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func WriteSuccessPaginated(w http.ResponseWriter, r *http.Request, status int, data any, meta PaginationMeta) {
	reqID := r.Header.Get("X-Request-ID")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success:   true,
		Data:      data,
		Meta:      meta,
		RequestId: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
