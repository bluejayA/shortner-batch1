package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/bluejayA/shortner-batch1/internal/handler"
	"github.com/bluejayA/shortner-batch1/internal/model"
)

// mockAuthServiceForHandler는 핸들러 테스트용 AuthService 목이다.
type mockAuthServiceForHandler struct {
	issueFn  func(ctx context.Context) (*model.APIKey, string, error)
	revokeFn func(ctx context.Context, key string) error
}

func (m *mockAuthServiceForHandler) Issue(ctx context.Context) (*model.APIKey, string, error) {
	return m.issueFn(ctx)
}
func (m *mockAuthServiceForHandler) Revoke(ctx context.Context, key string) error {
	return m.revokeFn(ctx, key)
}
func (m *mockAuthServiceForHandler) Validate(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func TestAuthHandler_Issue(t *testing.T) {
	svc := &mockAuthServiceForHandler{
		issueFn: func(_ context.Context) (*model.APIKey, string, error) {
			return &model.APIKey{ID: 1}, "plain-api-key", nil
		},
	}
	h := handler.NewAuthHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/keys", nil)
	w := httptest.NewRecorder()
	h.Issue(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("기대값 201, 실제값 %d", w.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("응답 JSON 파싱 실패: %v", err)
	}
	if resp["key"] != "plain-api-key" {
		t.Errorf("반환된 키가 다름: %s", resp["key"])
	}
}

func TestAuthHandler_Revoke(t *testing.T) {
	var revokedKey string
	svc := &mockAuthServiceForHandler{
		revokeFn: func(_ context.Context, key string) error {
			revokedKey = key
			return nil
		},
	}
	h := handler.NewAuthHandler(svc)

	r := chi.NewRouter()
	r.Delete("/api/keys/{key}", h.Revoke)

	req := httptest.NewRequest(http.MethodDelete, "/api/keys/test-key-123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("기대값 204, 실제값 %d", w.Code)
	}
	if revokedKey != "test-key-123" {
		t.Errorf("폐기된 키가 다름: %s", revokedKey)
	}
}
