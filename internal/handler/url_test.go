package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jay-ahn/shortner/internal/handler"
	"github.com/jay-ahn/shortner/internal/model"
	"github.com/jay-ahn/shortner/internal/service"
)

// mockURLService는 핸들러 테스트용 URLService 목이다.
type mockURLService struct {
	createFn  func(ctx context.Context, original, alias string, expiresAt *time.Time) (*model.URL, error)
	deleteFn  func(ctx context.Context, slug string) error
	resolveFn func(ctx context.Context, slug string) (string, error)
}

func (m *mockURLService) Create(ctx context.Context, original, alias string, expiresAt *time.Time) (*model.URL, error) {
	return m.createFn(ctx, original, alias, expiresAt)
}
func (m *mockURLService) Delete(ctx context.Context, slug string) error {
	return m.deleteFn(ctx, slug)
}
func (m *mockURLService) Resolve(ctx context.Context, slug string) (string, error) {
	return m.resolveFn(ctx, slug)
}

func TestURLHandler_Create(t *testing.T) {
	svc := &mockURLService{
		createFn: func(_ context.Context, _, _ string, _ *time.Time) (*model.URL, error) {
			return &model.URL{Slug: "abc123", Original: "https://example.com"}, nil
		},
	}
	h := handler.NewURLHandler(svc)

	body := `{"url":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("기대값 201, 실제값 %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["slug"] != "abc123" {
		t.Errorf("slug가 다름: %s", resp["slug"])
	}
}

func TestURLHandler_Delete(t *testing.T) {
	var deletedSlug string
	svc := &mockURLService{
		deleteFn: func(_ context.Context, slug string) error {
			deletedSlug = slug
			return nil
		},
	}
	h := handler.NewURLHandler(svc)

	r := chi.NewRouter()
	r.Delete("/api/urls/{slug}", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/urls/abc123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("기대값 204, 실제값 %d", w.Code)
	}
	if deletedSlug != "abc123" {
		t.Errorf("삭제된 slug가 다름: %s", deletedSlug)
	}
}

func TestRedirectHandler(t *testing.T) {
	r := chi.NewRouter()

	t.Run("유효한 slug는 302 리다이렉트", func(t *testing.T) {
		svc := &mockURLService{
			resolveFn: func(_ context.Context, _ string) (string, error) {
				return "https://example.com", nil
			},
		}
		h := handler.NewRedirectHandler(svc)
		r.Get("/{slug}", h.Redirect)

		req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("기대값 302, 실제값 %d", w.Code)
		}
		if w.Header().Get("Location") != "https://example.com" {
			t.Errorf("Location 헤더가 다름: %s", w.Header().Get("Location"))
		}
	})

	t.Run("없는 slug는 404", func(t *testing.T) {
		svc := &mockURLService{
			resolveFn: func(_ context.Context, _ string) (string, error) {
				return "", service.ErrNotFound
			},
		}
		h := handler.NewRedirectHandler(svc)
		r2 := chi.NewRouter()
		r2.Get("/{slug}", h.Redirect)

		req := httptest.NewRequest(http.MethodGet, "/notexist", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("기대값 404, 실제값 %d", w.Code)
		}
	})

	t.Run("만료된 slug는 410", func(t *testing.T) {
		svc := &mockURLService{
			resolveFn: func(_ context.Context, _ string) (string, error) {
				return "", service.ErrExpired
			},
		}
		h := handler.NewRedirectHandler(svc)
		r3 := chi.NewRouter()
		r3.Get("/{slug}", h.Redirect)

		req := httptest.NewRequest(http.MethodGet, "/expired", nil)
		w := httptest.NewRecorder()
		r3.ServeHTTP(w, req)

		if w.Code != http.StatusGone {
			t.Errorf("기대값 410, 실제값 %d", w.Code)
		}
	})
}

// errors 패키지 import 확인용
var _ = errors.New
