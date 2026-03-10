package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jay-ahn/shortner/internal/service"
)

// StatsHandler는 클릭 통계 조회 HTTP 핸들러다.
type StatsHandler struct {
	svc service.StatsService
}

// NewStatsHandler는 StatsHandler를 생성한다.
func NewStatsHandler(svc service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// Get은 GET /api/stats/{slug} — slug별 클릭 통계를 반환한다.
func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	stats, err := h.svc.Get(r.Context(), slug)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"slug":        stats.Slug,
		"click_count": stats.ClickCount,
	})
}
