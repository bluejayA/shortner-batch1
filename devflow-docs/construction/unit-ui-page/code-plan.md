# Code Generation Plan: unit-ui-page

## Test Files (Write First — TDD RED)
- [x] `internal/handler/ui_test.go` — GET / → 200, HTML 콘텐츠 포함 확인

## Test Strategy
- [x] `TestUIHandler_Get`: GET / → 200, Content-Type: text/html
- [x] `go test ./internal/handler/...` → GREEN 확인

## Implementation Files
- [x] `internal/static/embed.go` — index.html embed 패키지
- [x] `internal/static/index.html` — UI HTML (URL 단축 + 통계 조회)
- [x] `internal/handler/ui.go` — UIHandler (embed 서빙)

## Files to Modify
- [x] `internal/server/router.go` — GET / 라우트 추가 (variadic UIHandler)
- [x] `cmd/server/main.go` — UIHandler 초기화 + 라우터에 전달

## Implementation Steps
- [x] Step 1: 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: `web/index.html` 작성 (vanilla JS)
- [x] Step 4: `internal/static/embed.go` + `ui.go` 구현
- [x] Step 5: router에 GET / 추가
- [x] Step 6: `go test` 실행 → GREEN 확인
- [x] Step 7: 브라우저 동작 확인 → HTTP 200
