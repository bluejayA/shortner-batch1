# Code Generation Plan: unit-03-url

## Test Files (Write First — TDD RED)
- [x] `internal/cache/redis_test.go` — Cache 인터페이스 + ErrCacheMiss 검증
- [x] `internal/service/url_test.go` — URLService 비즈니스 로직 검증 (mock repo + cache)
- [x] `internal/handler/url_test.go` — URLHandler + RedirectHandler HTTP 응답 검증

## Test Strategy
- [x] `TestURLService_Create`: 자동 slug 생성, DB 저장 확인
- [x] `TestURLService_Create_CustomAlias`: 커스텀 alias가 slug로 사용됨
- [x] `TestURLService_Delete`: slug 삭제 + 캐시 무효화
- [x] `TestURLService_Resolve_CacheHit`: Redis에 slug 있으면 DB 조회 없이 반환
- [x] `TestURLService_Resolve_CacheMiss`: Redis 미스 시 DB 조회 + 캐시 저장
- [x] `TestURLService_Resolve_Expired`: 만료된 URL → ErrExpired 반환
- [x] `TestRedirectHandler`: 유효 slug → 302 / 없는 slug → 404 / 만료 → 410
- [x] `TestURLHandler_Create`: POST /api/urls → 201 + slug 반환
- [x] `TestURLHandler_Delete`: DELETE /api/urls/{slug} → 204
- [x] `go test ./internal/...` → GREEN 확인

## Implementation Files
- [x] `internal/cache/cache.go` — Cache 인터페이스 + Redis 구현체
- [x] `internal/repository/url.go` — PostgreSQL URLRepository 구현
- [x] `internal/service/url.go` — URLService 구현 (slug 생성 포함)
- [x] `internal/handler/redirect.go` — RedirectHandler 구현
- [x] `internal/handler/url.go` — URLHandler 구현

## Files to Modify
- [x] `internal/service/interfaces.go` — URLService 인터페이스 추가

## Implementation Steps
- [x] Step 1: 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: Cache 인터페이스 + Redis 구현
- [x] Step 4: URLRepository 구현
- [x] Step 5: URLService 구현
- [x] Step 6: RedirectHandler + URLHandler 구현
- [x] Step 7: `go test` 실행 → GREEN 확인 (16/16 PASS)
