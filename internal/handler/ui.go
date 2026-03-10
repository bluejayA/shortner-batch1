package handler

import (
	"net/http"

	"github.com/jay-ahn/shortner/internal/static"
)

// UIHandler는 웹 UI HTML 페이지를 서빙하는 핸들러다.
type UIHandler struct{}

// NewUIHandler는 UIHandler를 생성한다.
func NewUIHandler() *UIHandler {
	return &UIHandler{}
}

// Index는 GET / — index.html을 반환한다.
func (h *UIHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(static.IndexHTML)
}
