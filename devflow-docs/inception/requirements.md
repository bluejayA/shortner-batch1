# Requirements Analysis

**Depth**: Standard
**Timestamp**: 2026-03-10T00:00:00Z

## User Intent

Go 언어 기반의 URL 단축 서비스. 단축 URL 생성·리다이렉트 기본 기능 외에 커스텀 alias, 만료일, 클릭 통계, API 키 인증을 지원하는 프로덕션 수준의 백엔드 API.

## Functional Requirements

### URL 관리
- FR-01: 원본 URL을 입력받아 단축 코드(slug)를 자동 생성하여 반환
- FR-02: 사용자 지정 커스텀 alias 설정 가능 (slug 직접 지정)
- FR-03: 만료일(expiration) 설정 가능 — 만료 후 리다이렉트 시 410 Gone 반환
- FR-04: 단축 URL 삭제

### 리다이렉트
- FR-05: `GET /{slug}` 요청 시 원본 URL로 301/302 리다이렉트
- FR-06: 존재하지 않는 slug 요청 시 404 반환
- FR-07: 만료된 slug 요청 시 410 반환

### 통계
- FR-08: 리다이렉트 발생 시 클릭 수 기록
- FR-09: slug별 클릭 수 조회 API 제공

### 인증
- FR-10: API 키 기반 인증 — URL 생성/삭제/통계 조회는 유효한 API 키 필요
- FR-11: API 키 발급/폐기 API 제공
- FR-12: 리다이렉트(`GET /{slug}`)는 인증 불필요 (공개)

## Non-Functional Requirements

- NFR-01: 리다이렉트 응답시간 < 50ms (Redis 캐시 활용)
- NFR-02: API 응답시간 < 200ms (p99)
- NFR-03: slug 자동 생성 시 충돌 없는 고유값 보장
- NFR-04: API 키는 해시 저장 (평문 저장 금지)
- NFR-05: SQL Injection, 잘못된 URL 입력 방어

## Tech Stack

| 구성 요소 | 선택 |
|-----------|------|
| 언어 | Go |
| 웹 프레임워크 | net/http 표준 라이브러리 또는 chi |
| DB | PostgreSQL |
| 캐시 | Redis |
| 컨테이너 | Docker (멀티스테이지 빌드) |

## Assumptions

- 사용자 계정 시스템 없음 — API 키 단위로만 관리
- slug 길이: 6~8자 alphanumeric 자동 생성
- 통계는 단순 클릭 수만 우선 제공 (UA, 국가 등 상세 통계는 향후)
- HTTPS 종단처리는 인프라(nginx/LB) 레벨에서 담당

## Open Questions

- 없음 (모든 핵심 방향 결정됨)
