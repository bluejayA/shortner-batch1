# Code Generation Plan: unit-01-foundation

## Test Files (Write First — TDD RED)
- [x] `internal/model/url_test.go` — URL 만료 판별 로직 검증
- [x] `internal/model/apikey_test.go` — APIKey 해시 생성 로직 검증

## Test Strategy
- [x] `TestURL_IsExpired`: 만료일이 지난 URL은 true, 미설정 URL은 false 반환
- [x] `TestAPIKey_HashKey`: 동일 키 입력 시 동일 해시 반환, 빈 키 입력 시 에러
- [x] `go test ./internal/model/...` → RED 확인

## Implementation Files
- [x] `go.mod` + `go.sum` — 모듈 초기화 (chi, go-redis, lib/pq)
- [x] `internal/model/url.go` — URL 구조체 + `IsExpired()` 메서드
- [x] `internal/model/apikey.go` — APIKey 구조체 + `HashKey()` 유틸
- [x] `internal/model/stats.go` — Stats 구조체
- [x] `migrations/001_init.sql` — urls, api_keys, click_stats 테이블 DDL
- [x] `docker-compose.yml` — PostgreSQL 5432, Redis 6379
- [x] `Dockerfile` — 멀티스테이지 빌드 (builder + runner)
- [x] `.env.example` — DB/Redis 접속 정보 템플릿
- [x] `.gitignore` — .env, 바이너리 등 제외

## Implementation Steps
- [x] Step 1: 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: go.mod 초기화 + 의존성 추가
- [x] Step 4: 모델 파일 구현
- [x] Step 5: `go test` 실행 → GREEN 확인 (6/6 PASS)
- [x] Step 6: 마이그레이션 SQL + Docker 파일 작성
