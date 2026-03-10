package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bluejayA/shortner-batch1/internal/handler"
)

func TestUIHandler_Get(t *testing.T) {
	h := handler.NewUIHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("기대값 200, 실제값 %d", w.Code)
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Errorf("Content-Type이 text/html이 아님: %s", w.Header().Get("Content-Type"))
	}
	if !strings.Contains(w.Body.String(), "<html") {
		t.Error("응답 바디에 HTML이 없음")
	}
}
