# Code Generation Plan: unit-04-stats

## Test Files (Write First — TDD RED)
- [x] `internal/service/stats_test.go` — StatsService Record/Get 검증 (mock repo)
- [x] `internal/handler/stats_test.go` — StatsHandler HTTP 응답 검증

## Test Strategy
- [x] `TestStatsService_Record`: slug 클릭 수 증가 호출 확인
- [x] `TestStatsService_Get`: slug별 통계 반환 확인
- [x] `TestStatsHandler_Get`: GET /api/stats/{slug} → 200 + 클릭 수 JSON
- [x] `TestStatsHandler_Get_NotFound`: 없는 slug → 404
- [x] `go test ./internal/...` → GREEN 확인

## Implementation Files
- [x] `internal/repository/stats.go` — PostgreSQL StatsRepository 구현
- [x] `internal/service/stats.go` — StatsService 구현
- [x] `internal/handler/stats.go` — StatsHandler 구현

## Files to Modify
- [x] `internal/service/interfaces.go` — StatsService 인터페이스 추가
- [x] `internal/handler/redirect.go` — 비동기 goroutine으로 StatsService.Record 호출 추가

## Implementation Steps
- [x] Step 1: 테스트 파일 작성 (RED)
- [x] Step 2: `go test` 실행 → RED 확인
- [x] Step 3: StatsRepository 구현
- [x] Step 4: StatsService 구현
- [x] Step 5: StatsHandler 구현
- [x] Step 6: RedirectHandler에 통계 기록 통합
- [x] Step 7: `go test` 실행 → GREEN 확인 (22/22 PASS)
