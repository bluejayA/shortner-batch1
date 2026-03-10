package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jay-ahn/shortner/internal/service"
)

// AuthHandler는 API 키 관련 HTTP 핸들러다.
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler는 AuthHandler를 생성한다.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Issue는 POST /api/keys — 새 API 키를 발급한다.
func (h *AuthHandler) Issue(w http.ResponseWriter, r *http.Request) {
	_, plaintext, err := h.svc.Issue(r.Context())
	if err != nil {
		http.Error(w, "failed to issue key", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"key": plaintext})
}

// Revoke는 DELETE /api/keys/{key} — API 키를 폐기한다.
func (h *AuthHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if err := h.svc.Revoke(r.Context(), key); err != nil {
		http.Error(w, "failed to revoke key", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
