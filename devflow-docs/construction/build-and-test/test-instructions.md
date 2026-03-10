# Test Instructions

## Unit Tests

모든 테스트는 외부 의존성(DB, Redis) 없이 mock으로 실행됩니다.

```bash
go test ./... -v
```

**Expected**: 43개 테스트 PASS, 6개 패키지 전체 통과

```
ok  github.com/jay-ahn/shortner/internal/cache
ok  github.com/jay-ahn/shortner/internal/handler
ok  github.com/jay-ahn/shortner/internal/middleware
ok  github.com/jay-ahn/shortner/internal/model
ok  github.com/jay-ahn/shortner/internal/server
ok  github.com/jay-ahn/shortner/internal/service
```

패키지별 개별 실행:

```bash
go test ./internal/model/...      # 모델 테스트 (3개)
go test ./internal/service/...    # 서비스 로직 (14개)
go test ./internal/middleware/...  # 인증 미들웨어 (3개)
go test ./internal/handler/...    # HTTP 핸들러 (9개)
go test ./internal/server/...     # 라우터 통합 (4개)
go test ./internal/cache/...      # 캐시 인터페이스 (1개)
```

## Manual Verification (통합 테스트)

로컬 환경 기동 후 아래 순서로 검증합니다.

### 1. 서버 기동 확인

```bash
curl http://localhost:8080/health
# Expected: 200 ok
```

### 2. API 키 발급

```bash
curl -X POST http://localhost:8080/api/keys
# Expected: 201 {"key": "<plaintext-key>"}
```

`<plaintext-key>` 값을 이후 요청에 사용합니다.

### 3. URL 단축 생성

```bash
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer <plaintext-key>" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
# Expected: 201 {"slug": "abc123", "short_url": "/abc123"}
```

### 4. 리다이렉트 확인 (공개, 인증 불필요)

```bash
curl -L http://localhost:8080/<slug>
# Expected: 302 → https://www.google.com 으로 리다이렉트
```

### 5. 클릭 통계 확인

```bash
curl http://localhost:8080/api/stats/<slug> \
  -H "Authorization: Bearer <plaintext-key>"
# Expected: 200 {"slug": "...", "click_count": 1}
```

### 6. 커스텀 alias + 만료일

```bash
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer <plaintext-key>" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "alias": "my-link", "expires_at": "2026-12-31T00:00:00Z"}'
# Expected: 201 {"slug": "my-link", ...}
```

### 7. 인증 없는 보호 라우트 → 401

```bash
curl -X POST http://localhost:8080/api/urls \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
# Expected: 401 unauthorized
```

## Coverage Gaps

- **Repository 레이어**: 실제 DB 연동 통합 테스트 없음 (mock으로 대체)
  - 향후 `testcontainers-go`를 활용한 PostgreSQL 통합 테스트 추가 권장
- **Redis 캐시**: 실제 Redis 연동 테스트 없음
