package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/bluejayA/shortner-batch1/internal/service"
)

// URLHandler는 URL 생성·삭제 HTTP 핸들러다.
type URLHandler struct {
	svc service.URLService
}

// NewURLHandler는 URLHandler를 생성한다.
func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

// createRequest는 URL 생성 요청 바디다.
type createRequest struct {
	URL       string     `json:"url"`
	Alias     string     `json:"alias,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Create는 POST /api/urls — 새 단축 URL을 생성한다.
func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.svc.Create(r.Context(), req.URL, req.Alias, req.ExpiresAt)
	if err != nil {
		http.Error(w, "failed to create url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"slug":     url.Slug,
		"short_url": "/" + url.Slug,
	})
}

// Delete는 DELETE /api/urls/{slug} — 단축 URL을 삭제한다.
func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if err := h.svc.Delete(r.Context(), slug); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
