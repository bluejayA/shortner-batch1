package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/bluejayA/shortner-batch1/internal/handler"
	"github.com/bluejayA/shortner-batch1/internal/middleware"
)

// NewRouter는 모든 핸들러와 미들웨어를 조립한 chi 라우터를 반환한다.
func NewRouter(
	redirectH *handler.RedirectHandler,
	urlH *handler.URLHandler,
	authH *handler.AuthHandler,
	statsH *handler.StatsHandler,
	authMw *middleware.AuthMiddleware,
	uiH ...*handler.UIHandler,
) http.Handler {
	r := chi.NewRouter()

	// 공통 미들웨어
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	// UI (공개)
	if len(uiH) > 0 && uiH[0] != nil {
		r.Get("/", uiH[0].Index)
	}

	// 헬스체크 (공개)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// 리다이렉트 (공개)
	r.Get("/{slug}", redirectH.Redirect)

	// API 키 발급 (공개 — 최초 키 발급을 위해 인증 불필요)
	r.Post("/api/keys", authH.Issue)

	// 인증 필요 라우트
	r.Group(func(r chi.Router) {
		r.Use(authMw.Middleware)

		// URL 관리
		r.Post("/api/urls", urlH.Create)
		r.Delete("/api/urls/{slug}", urlH.Delete)

		// 통계
		r.Get("/api/stats/{slug}", statsH.Get)

		// API 키 폐기 (인증 필요)
		r.Delete("/api/keys/{key}", authH.Revoke)
	})

	return r
}
