# Units Generation

**Depth**: Minimal
**Timestamp**: 2026-03-10T00:04:00Z

---

### Unit: unit-01-foundation
**Responsibility**: Go 모듈 초기화, 도메인 모델 정의, DB 마이그레이션, Docker 환경 구성.
**Dependencies**: 없음
**Interfaces**:
- `internal/model/` — `URL`, `APIKey`, `Stats` 구조체
- `migrations/` — PostgreSQL DDL SQL 파일
- `docker-compose.yml` — PostgreSQL + Redis 로컬 실행 환경
- `go.mod` — chi, go-redis, lib/pq 의존성
**Implementation order**: 1

---

### Unit: unit-02-auth
**Responsibility**: API 키 발급·폐기·검증 (APIKeyRepository + AuthService + AuthMiddleware + AuthHandler).
**Dependencies**: unit-01-foundation (model.APIKey)
**Interfaces**:
- `AuthService.Issue(ctx) (*model.APIKey, error)`
- `AuthService.Revoke(ctx, key string) error`
- `AuthService.Validate(ctx, key string) (bool, error)`
- `AuthMiddleware.Middleware(next http.Handler) http.Handler`
- HTTP: `POST /api/keys`, `DELETE /api/keys/{key}`
**Implementation order**: 2 (unit-03, unit-04와 병렬 가능)

---

### Unit: unit-03-url
**Responsibility**: URL 단축 핵심 기능 — 생성·삭제·리다이렉트 (URLRepository + Cache + URLService + URLHandler + RedirectHandler).
**Dependencies**: unit-01-foundation (model.URL)
**Interfaces**:
- `URLService.Create(ctx, originalURL, alias string, expiresAt *time.Time) (*model.URL, error)`
- `URLService.Delete(ctx, slug string) error`
- `URLService.Resolve(ctx, slug string) (string, error)`
- `Cache.Get/Set/Delete(ctx, slug string, ...)`
- HTTP: `GET /{slug}`, `POST /api/urls`, `DELETE /api/urls/{slug}`
**Implementation order**: 2 (unit-02, unit-04와 병렬 가능)

---

### Unit: unit-04-stats
**Responsibility**: 클릭 통계 기록 및 조회 (StatsRepository + StatsService + StatsHandler).
**Dependencies**: unit-01-foundation (model.Stats)
**Interfaces**:
- `StatsService.Record(ctx, slug string) error`
- `StatsService.Get(ctx, slug string) (*model.Stats, error)`
- HTTP: `GET /api/stats/{slug}`
**Implementation order**: 2 (unit-02, unit-03와 병렬 가능)

---

### Unit: unit-05-server
**Responsibility**: Router 구성, 미들웨어 체인 조립, main.go 진입점, 전체 통합 테스트.
**Dependencies**: unit-02-auth, unit-03-url, unit-04-stats (모두 완료 후)
**Interfaces**:
- `cmd/server/main.go` — 서버 기동 진입점
- 전체 라우트 테이블 조립
**Implementation order**: 3
