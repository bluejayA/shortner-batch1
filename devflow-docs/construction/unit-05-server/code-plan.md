# Code Generation Plan: unit-05-server

## Test Files (Write First — TDD RED)
- [x] `internal/server/router_test.go` — 라우터 통합: 라우트 등록 확인, 미들웨어 동작 검증

## Test Strategy
- [x] `TestRouter_HealthCheck`: GET /health → 200
- [x] `TestRouter_RedirectPublic`: GET /{slug} 인증 없이 접근 가능 (공개 라우트)
- [x] `TestRouter_ProtectedRoute_NoAuth`: POST /api/urls → 인증 없이 401
- [x] `TestRouter_ProtectedRoute_WithAuth`: POST /api/urls + 유효 토큰 → 201
- [x] `go test ./internal/server/...` → GREEN 확인

## Implementation Files
- [x] `internal/server/router.go` — chi 라우터 조립 (핸들러 + 미들웨어 연결)
- [x] `cmd/server/main.go` — DB/Redis 연결, 서비스 초기화, 서버 기동

## Implementation Steps
- [x] Step 1: 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: `internal/server/router.go` 구현
- [x] Step 4: `go test` 실행 → GREEN 확인 (4/4 PASS)
- [x] Step 5: `cmd/server/main.go` 구현
- [x] Step 6: `go build ./cmd/server` 빌드 성공 확인
