package middleware

import (
	"net/http"
	"strings"

	"github.com/bluejayA/shortner-batch1/internal/service"
)

// AuthMiddleware는 API 키 인증 미들웨어다.
type AuthMiddleware struct {
	svc service.AuthService
}

// NewAuthMiddleware는 AuthMiddleware를 생성한다.
func NewAuthMiddleware(svc service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{svc: svc}
}

// Middleware는 Authorization: Bearer <key> 헤더를 검증한다.
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		key := strings.TrimPrefix(authHeader, "Bearer ")
		ok, err := m.svc.Validate(r.Context(), key)
		if err != nil || !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
