package server_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jay-ahn/shortner/internal/handler"
	"github.com/jay-ahn/shortner/internal/middleware"
	"github.com/jay-ahn/shortner/internal/model"
	"github.com/jay-ahn/shortner/internal/server"
	"github.com/jay-ahn/shortner/internal/service"
)

// --- mock 서비스 ---

type mockURLSvc struct{}

func (m *mockURLSvc) Create(_ context.Context, _, _ string, _ *time.Time) (*model.URL, error) {
	return &model.URL{Slug: "abc123", Original: "https://example.com"}, nil
}
func (m *mockURLSvc) Delete(_ context.Context, _ string) error { return nil }
func (m *mockURLSvc) Resolve(_ context.Context, slug string) (string, error) {
	if slug == "abc123" {
		return "https://example.com", nil
	}
	return "", service.ErrNotFound
}

type mockAuthSvc struct{}

func (m *mockAuthSvc) Issue(_ context.Context) (*model.APIKey, string, error) {
	return &model.APIKey{ID: 1}, "test-key", nil
}
func (m *mockAuthSvc) Revoke(_ context.Context, _ string) error { return nil }
func (m *mockAuthSvc) Validate(_ context.Context, key string) (bool, error) {
	return key == "valid-key", nil
}

type mockStatsSvc struct{}

func (m *mockStatsSvc) Record(_ context.Context, _ string) error          { return nil }
func (m *mockStatsSvc) Get(_ context.Context, slug string) (*model.Stats, error) {
	return &model.Stats{Slug: slug, ClickCount: 0}, nil
}

// --- 테스트 헬퍼: 테스트용 라우터 생성 ---
func newTestRouter() http.Handler {
	urlSvc := &mockURLSvc{}
	authSvc := &mockAuthSvc{}
	statsSvc := &mockStatsSvc{}

	redirectH := handler.NewRedirectHandler(urlSvc, statsSvc)
	urlH := handler.NewURLHandler(urlSvc)
	authH := handler.NewAuthHandler(authSvc)
	statsH := handler.NewStatsHandler(statsSvc)
	authMw := middleware.NewAuthMiddleware(authSvc)

	return server.NewRouter(redirectH, urlH, authH, statsH, authMw)
}

func TestRouter_HealthCheck(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("기대값 200, 실제값 %d", w.Code)
	}
}

func TestRouter_RedirectPublic(t *testing.T) {
	r := newTestRouter()
	// 인증 없이 리다이렉트 라우트 접근 가능해야 함
	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("공개 리다이렉트 기대값 302, 실제값 %d", w.Code)
	}
}

func TestRouter_ProtectedRoute_NoAuth(t *testing.T) {
	r := newTestRouter()
	// 인증 없이 보호된 라우트 → 401
	body := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("인증 없는 요청 기대값 401, 실제값 %d", w.Code)
	}
}

func TestRouter_ProtectedRoute_WithAuth(t *testing.T) {
	r := newTestRouter()
	// 유효한 API 키로 보호된 라우트 → 201
	body := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("인증된 요청 기대값 201, 실제값 %d", w.Code)
	}
}
