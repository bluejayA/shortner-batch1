# Application Design

**Depth**: Standard
**Timestamp**: 2026-03-10T00:03:00Z

---

## 디렉토리 구조

```
shortner/
├── cmd/
│   └── server/
│       └── main.go              # 진입점
├── internal/
│   ├── handler/                 # HTTP 핸들러 (라우팅 레이어)
│   │   ├── redirect.go
│   │   ├── url.go
│   │   ├── auth.go
│   │   └── stats.go
│   ├── service/                 # 비즈니스 로직 레이어
│   │   ├── url.go
│   │   ├── auth.go
│   │   └── stats.go
│   ├── repository/              # 데이터 접근 레이어 (PostgreSQL)
│   │   ├── url.go
│   │   ├── apikey.go
│   │   └── stats.go
│   ├── middleware/              # HTTP 미들웨어
│   │   └── auth.go
│   ├── model/                   # 도메인 모델
│   │   ├── url.go
│   │   ├── apikey.go
│   │   └── stats.go
│   └── cache/                   # Redis 캐시 래퍼
│       └── redis.go
├── migrations/                  # DB 마이그레이션 SQL
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

---

## 컴포넌트 설계

### Component: Router
- **책임**: HTTP 요청을 핸들러로 라우팅, 미들웨어 체인 구성
- **Public Interface**:
  - `GET /{slug}` → RedirectHandler (공개)
  - `POST /api/urls` → URLHandler (인증 필요)
  - `DELETE /api/urls/{slug}` → URLHandler (인증 필요)
  - `GET /api/stats/{slug}` → StatsHandler (인증 필요)
  - `POST /api/keys` → AuthHandler (인증 필요)
  - `DELETE /api/keys/{key}` → AuthHandler (인증 필요)
- **Dependencies**: chi router
- **Data Owned**: 없음

---

### Component: RedirectHandler
- **책임**: `GET /{slug}` 요청을 받아 원본 URL로 리다이렉트
- **Public Interface**: `ServeHTTP(w, r)`
- **Dependencies**: URLService, Cache
- **Data Owned**: 없음

---

### Component: URLHandler
- **책임**: URL 생성·삭제 HTTP 요청 처리
- **Public Interface**:
  - `Create(w, r)` — POST /api/urls
  - `Delete(w, r)` — DELETE /api/urls/{slug}
- **Dependencies**: URLService
- **Data Owned**: 없음

---

### Component: AuthHandler
- **책임**: API 키 발급·폐기 HTTP 요청 처리
- **Public Interface**:
  - `Issue(w, r)` — POST /api/keys
  - `Revoke(w, r)` — DELETE /api/keys/{key}
- **Dependencies**: AuthService
- **Data Owned**: 없음

---

### Component: StatsHandler
- **책임**: 클릭 통계 조회 HTTP 요청 처리
- **Public Interface**:
  - `Get(w, r)` — GET /api/stats/{slug}
- **Dependencies**: StatsService
- **Data Owned**: 없음

---

### Component: AuthMiddleware
- **책임**: 보호된 라우트에서 API 키 유효성 검증
- **Public Interface**: `Middleware(next http.Handler) http.Handler`
- **Dependencies**: AuthService
- **Data Owned**: 없음
- **동작**:
  - `Authorization: Bearer <key>` 헤더 추출
  - 유효하지 않으면 401 반환
  - 유효하면 context에 key 정보 주입 후 next 호출

---

### Component: URLService
- **책임**: URL 단축 비즈니스 로직 (slug 생성, 중복 확인, 만료 처리)
- **Public Interface**:
  - `Create(ctx, originalURL, alias, expiresAt) (*model.URL, error)`
  - `Delete(ctx, slug) error`
  - `Resolve(ctx, slug) (string, error)`  — 원본 URL 반환, 만료 시 ErrExpired
- **Dependencies**: URLRepository, Cache
- **Data Owned**: 없음
- **slug 생성 전략**: 6자 alphanumeric 랜덤 생성 → 충돌 시 재시도 (최대 3회)

---

### Component: AuthService
- **책임**: API 키 발급·폐기·검증 비즈니스 로직
- **Public Interface**:
  - `Issue(ctx) (*model.APIKey, error)` — 새 키 생성 및 해시 저장
  - `Revoke(ctx, key) error`
  - `Validate(ctx, key) (bool, error)` — 해시 비교로 검증
- **Dependencies**: APIKeyRepository
- **Data Owned**: 없음
- **보안**: bcrypt 또는 SHA-256으로 키 해시 저장, 평문 반환은 발급 시 1회만

---

### Component: StatsService
- **책임**: 클릭 이벤트 기록 및 조회
- **Public Interface**:
  - `Record(ctx, slug) error`
  - `Get(ctx, slug) (*model.Stats, error)`
- **Dependencies**: StatsRepository
- **Data Owned**: 없음

---

### Component: URLRepository
- **책임**: PostgreSQL URLs 테이블 CRUD
- **Public Interface**:
  - `Insert(ctx, url *model.URL) error`
  - `FindBySlug(ctx, slug) (*model.URL, error)`
  - `DeleteBySlug(ctx, slug) error`
- **Dependencies**: PostgreSQL 커넥션 풀
- **Data Owned**: `urls` 테이블

---

### Component: APIKeyRepository
- **책임**: PostgreSQL API 키 테이블 CRUD
- **Public Interface**:
  - `Insert(ctx, key *model.APIKey) error`
  - `FindByHash(ctx, hash) (*model.APIKey, error)`
  - `DeleteByHash(ctx, hash) error`
- **Dependencies**: PostgreSQL 커넥션 풀
- **Data Owned**: `api_keys` 테이블

---

### Component: StatsRepository
- **책임**: PostgreSQL 클릭 통계 CRUD
- **Public Interface**:
  - `Increment(ctx, slug) error`
  - `FindBySlug(ctx, slug) (*model.Stats, error)`
- **Dependencies**: PostgreSQL 커넥션 풀
- **Data Owned**: `click_stats` 테이블

---

### Component: Cache
- **책임**: Redis 기반 slug → 원본 URL 캐싱 (리다이렉트 성능 최적화)
- **Public Interface**:
  - `Get(ctx, slug) (string, error)`
  - `Set(ctx, slug, url string, ttl time.Duration) error`
  - `Delete(ctx, slug) error`
- **Dependencies**: Redis 클라이언트 (go-redis)
- **Data Owned**: Redis key-value

---

## DB 스키마

```sql
-- URLs 테이블
CREATE TABLE urls (
    slug        VARCHAR(16) PRIMARY KEY,
    original    TEXT NOT NULL,
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- API 키 테이블
CREATE TABLE api_keys (
    id          SERIAL PRIMARY KEY,
    key_hash    VARCHAR(64) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

-- 클릭 통계 테이블
CREATE TABLE click_stats (
    slug        VARCHAR(16) PRIMARY KEY REFERENCES urls(slug) ON DELETE CASCADE,
    click_count BIGINT DEFAULT 0,
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 컴포넌트 상호작용 다이어그램

```
# 리다이렉트 플로우 (공개, 고성능)
[Request GET /{slug}]
  --> [Router]
  --> [RedirectHandler]
  --> [Cache.Get(slug)]     # Redis 우선 조회
      Hit  --> 302 Redirect
      Miss --> [URLService.Resolve(slug)]
               --> [URLRepository.FindBySlug]  # PostgreSQL
               --> [Cache.Set(slug, url)]       # 캐시 워밍
               --> 302 Redirect (or 404/410)
  --> [StatsService.Record(slug)]               # 비동기 기록 (goroutine)

# URL 생성 플로우 (인증 필요)
[Request POST /api/urls]
  --> [Router]
  --> [AuthMiddleware] --> [AuthService.Validate] --> [APIKeyRepository]
  --> [URLHandler.Create]
  --> [URLService.Create]
  --> [URLRepository.Insert]
  --> 201 Created

# API 키 발급 플로우
[Request POST /api/keys]
  --> [Router]
  --> [AuthHandler.Issue]
  --> [AuthService.Issue]
  --> [APIKeyRepository.Insert]   # 해시만 저장
  --> 201 Created (평문 키 1회 반환)
```

---

## Design Decisions

### 1. 웹 프레임워크: chi vs net/http 표준
- **chi 선택**: URL 파라미터 추출(`/{slug}`)이 표준 라이브러리보다 간결, 미들웨어 체인 구성 편의성
- 외부 의존성 최소화 원칙 하에 chi는 매우 경량 (net/http 래퍼 수준)

### 2. 클릭 통계 기록 방식: 동기 vs 비동기
- **비동기(goroutine) 선택**: 리다이렉트 응답시간 < 50ms NFR 달성을 위해 통계 기록이 응답 지연에 영향을 주지 않도록 분리
- 단순 `go statsService.Record(...)` 패턴 사용 (현재 규모에서 메시지 큐 불필요)

### 3. API 키 해시 알고리즘: bcrypt vs SHA-256
- **SHA-256 선택**: bcrypt는 의도적으로 느린 알고리즘 → 매 요청마다 미들웨어에서 검증 시 성능 저하
- SHA-256으로 해시 후 DB 조회 방식 (검증 자체가 DB 조회 결과로 이루어짐)
