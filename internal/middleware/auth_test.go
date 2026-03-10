package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jay-ahn/shortner/internal/middleware"
	"github.com/jay-ahn/shortner/internal/model"
)

// mockAuthService는 테스트용 AuthService 목이다.
type mockAuthService struct {
	validateFn func(ctx context.Context, key string) (bool, error)
}

func (m *mockAuthService) Issue(_ context.Context) (*model.APIKey, string, error) {
	return nil, "", nil
}
func (m *mockAuthService) Revoke(_ context.Context, _ string) error { return nil }
func (m *mockAuthService) Validate(ctx context.Context, key string) (bool, error) {
	return m.validateFn(ctx, key)
}

func TestAuthMiddleware(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	t.Run("유효한 Bearer 토큰은 next를 호출한다", func(t *testing.T) {
		nextCalled = false
		svc := &mockAuthService{
			validateFn: func(_ context.Context, key string) (bool, error) {
				return key == "valid-key", nil
			},
		}
		mw := middleware.NewAuthMiddleware(svc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer valid-key")
		w := httptest.NewRecorder()
		mw.Middleware(next).ServeHTTP(w, req)
		if !nextCalled {
			t.Error("유효한 토큰인데 next가 호출되지 않음")
		}
		if w.Code != http.StatusOK {
			t.Errorf("기대값 200, 실제값 %d", w.Code)
		}
	})

	t.Run("Authorization 헤더 없으면 401 반환", func(t *testing.T) {
		nextCalled = false
		svc := &mockAuthService{
			validateFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
		}
		mw := middleware.NewAuthMiddleware(svc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mw.Middleware(next).ServeHTTP(w, req)
		if nextCalled {
			t.Error("헤더 없는데 next가 호출됨")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("기대값 401, 실제값 %d", w.Code)
		}
	})

	t.Run("잘못된 토큰은 401 반환", func(t *testing.T) {
		nextCalled = false
		svc := &mockAuthService{
			validateFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
		}
		mw := middleware.NewAuthMiddleware(svc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-key")
		w := httptest.NewRecorder()
		mw.Middleware(next).ServeHTTP(w, req)
		if nextCalled {
			t.Error("잘못된 토큰인데 next가 호출됨")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("기대값 401, 실제값 %d", w.Code)
		}
	})
}
