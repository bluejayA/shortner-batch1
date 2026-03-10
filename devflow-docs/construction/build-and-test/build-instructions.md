# Build Instructions

## Prerequisites

- Go 1.22+
- Docker & Docker Compose (로컬 개발 환경)

## 로컬 개발 환경 시작

```bash
# PostgreSQL + Redis 컨테이너 기동
docker compose up -d

# 컨테이너 상태 확인
docker compose ps
```

## 의존성 동기화

```bash
go mod tidy
```

## 바이너리 빌드

```bash
go build -o shortner ./cmd/server
```

## 실행

```bash
# .env 파일 생성 (최초 1회)
cp .env.example .env

# 환경변수 적용 후 실행
export $(cat .env | xargs) && ./shortner
```

또는 환경변수 직접 지정:

```bash
PORT=8080 \
DATABASE_URL=postgres://shortner:shortner@localhost:5432/shortner?sslmode=disable \
REDIS_URL=redis://localhost:6379/0 \
./shortner
```

## Expected Output

```
서버 시작: http://localhost:8080
```

## Docker 이미지 빌드 (멀티스테이지)

```bash
docker build -t shortner:latest .
docker run -p 8080:8080 \
  -e DATABASE_URL=... \
  -e REDIS_URL=... \
  shortner:latest
```
