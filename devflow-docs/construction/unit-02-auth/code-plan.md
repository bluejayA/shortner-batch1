# Code Generation Plan: unit-02-auth

## Test Files (Write First — TDD RED)
- [x] `internal/repository/apikey_test.go` — APIKeyRepository 인터페이스 정의 및 mock
- [x] `internal/service/auth_test.go` — AuthService 비즈니스 로직 검증 (mock repo)
- [x] `internal/middleware/auth_test.go` — AuthMiddleware 토큰 검증 동작 검증
- [x] `internal/handler/auth_test.go` — AuthHandler HTTP 응답 검증

## Test Strategy
- [x] `TestAuthService_Issue`: 키 발급 시 plaintext 반환 + hash 저장 확인
- [x] `TestAuthService_Revoke`: 유효한 키 폐기 성공 / 존재하지 않는 키 에러
- [x] `TestAuthService_Validate`: 유효한 키 true / 폐기된 키 false / 없는 키 false
- [x] `TestAuthMiddleware`: 유효한 Bearer 토큰 → next 호출 / 없거나 잘못된 토큰 → 401
- [x] `TestAuthHandler_Issue`: POST /api/keys → 201 + 키 반환
- [x] `TestAuthHandler_Revoke`: DELETE /api/keys/{key} → 204
- [x] `go test ./internal/...` → RED 확인

## Implementation Files
- [x] `internal/repository/interfaces.go` — 모든 repository 인터페이스 정의
- [x] `internal/repository/apikey.go` — PostgreSQL APIKeyRepository 구현
- [x] `internal/service/interfaces.go` — AuthService 인터페이스 정의
- [x] `internal/service/auth.go` — AuthService 구현
- [x] `internal/middleware/auth.go` — AuthMiddleware 구현
- [x] `internal/handler/auth.go` — AuthHandler 구현

## Implementation Steps
- [x] Step 1: 인터페이스 + 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: repository 인터페이스 + PostgreSQL 구현
- [x] Step 4: AuthService 구현
- [x] Step 5: AuthMiddleware 구현
- [x] Step 6: AuthHandler 구현
- [x] Step 7: `go test` 실행 → GREEN 확인 (11/11 PASS)
