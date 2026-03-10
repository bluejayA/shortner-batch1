package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/bluejayA/shortner-batch1/internal/handler"
	"github.com/bluejayA/shortner-batch1/internal/model"
)

// mockStatsService는 핸들러 테스트용 StatsService 목이다.
type mockStatsService struct {
	recordFn func(ctx context.Context, slug string) error
	getFn    func(ctx context.Context, slug string) (*model.Stats, error)
}

func (m *mockStatsService) Record(ctx context.Context, slug string) error {
	return m.recordFn(ctx, slug)
}
func (m *mockStatsService) Get(ctx context.Context, slug string) (*model.Stats, error) {
	return m.getFn(ctx, slug)
}

func TestStatsHandler_Get(t *testing.T) {
	svc := &mockStatsService{
		getFn: func(_ context.Context, slug string) (*model.Stats, error) {
			return &model.Stats{Slug: slug, ClickCount: 99}, nil
		},
	}
	h := handler.NewStatsHandler(svc)

	r := chi.NewRouter()
	r.Get("/api/stats/{slug}", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/stats/abc123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("기대값 200, 실제값 %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["slug"] != "abc123" {
		t.Errorf("slug가 다름: %v", resp["slug"])
	}
	if int(resp["click_count"].(float64)) != 99 {
		t.Errorf("클릭 수가 다름: %v", resp["click_count"])
	}
}

func TestStatsHandler_Get_NotFound(t *testing.T) {
	svc := &mockStatsService{
		getFn: func(_ context.Context, _ string) (*model.Stats, error) {
			return nil, errors.New("not found")
		},
	}
	h := handler.NewStatsHandler(svc)

	r := chi.NewRouter()
	r.Get("/api/stats/{slug}", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/stats/notexist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("기대값 404, 실제값 %d", w.Code)
	}
}
