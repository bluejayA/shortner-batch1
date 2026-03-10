package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jay-ahn/shortner/internal/service"
)

// RedirectHandler는 단축 URL 리다이렉트 핸들러다.
type RedirectHandler struct {
	urlSvc   service.URLService
	statsSvc service.StatsService
}

// NewRedirectHandler는 RedirectHandler를 생성한다.
func NewRedirectHandler(urlSvc service.URLService, statsSvc ...service.StatsService) *RedirectHandler {
	h := &RedirectHandler{urlSvc: urlSvc}
	if len(statsSvc) > 0 {
		h.statsSvc = statsSvc[0]
	}
	return h
}

// Redirect는 GET /{slug} — 원본 URL로 302 리다이렉트하고 클릭 수를 비동기 기록한다.
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	original, err := h.urlSvc.Resolve(r.Context(), slug)
	if err != nil {
		if errors.Is(err, service.ErrExpired) {
			http.Error(w, "gone", http.StatusGone)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// 클릭 통계 비동기 기록 — context.Background() 사용 (request context는 응답 후 취소됨)
	if h.statsSvc != nil {
		go func() {
			_ = h.statsSvc.Record(context.Background(), slug)
		}()
	}

	http.Redirect(w, r, original, http.StatusFound)
}
