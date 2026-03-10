# Shortner

> A fast, self-hosted URL shortener built with Go.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-43%20passing-brightgreen)](#)

---

## Features

- **URL Shortening** — Shorten any URL with auto-generated or custom slugs
- **Redirect** — Fast 302 redirects with Redis caching (< 50ms)
- **Expiration** — Set an expiry date on any short URL (returns `410 Gone` when expired)
- **Click Stats** — Track click counts per slug
- **API Key Auth** — Secure write operations with SHA-256 hashed API keys
- **Web UI** — Minimal browser interface for creating and looking up short URLs
- **Docker-ready** — Multi-stage Dockerfile + Docker Compose for local development

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22+ |
| Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Container | Docker (multi-stage build) |

---

## Quick Start

### Prerequisites

- [Go 1.22+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) & Docker Compose

### 1. Clone

```bash
git clone https://github.com/bluejayA/shortner-batch1.git
cd shortner-batch1
```

### 2. Start dependencies

```bash
docker compose up -d
```

### 3. Run the server

```bash
go run ./cmd/server
# Server listening on http://localhost:8080
```

### 4. Open the Web UI

Visit **http://localhost:8080** in your browser.

---

## API Reference

All write endpoints require an `Authorization: Bearer <api-key>` header.

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/keys` | ✗ | Issue a new API key |
| `DELETE` | `/api/keys/{key}` | ✓ | Revoke an API key |

#### Issue API key

```bash
curl -X POST http://localhost:8080/api/keys
```

```json
{ "key": "d45b83d8975f41b8..." }
```

> ⚠️ The plaintext key is returned **once only**. Store it securely.

---

### URLs

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/urls` | ✓ | Create a short URL |
| `DELETE` | `/api/urls/{slug}` | ✓ | Delete a short URL |
| `GET` | `/{slug}` | ✗ | Redirect to original URL |

#### Create a short URL

```bash
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer <your-key>" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/very/long/path",
    "alias": "my-link",
    "expires_at": "2026-12-31T00:00:00Z"
  }'
```

```json
{ "slug": "my-link", "short_url": "/my-link" }
```

| Field | Required | Description |
|-------|----------|-------------|
| `url` | ✓ | Original URL to shorten |
| `alias` | ✗ | Custom slug (auto-generated if omitted) |
| `expires_at` | ✗ | ISO 8601 expiry datetime |

---

### Stats

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/stats/{slug}` | ✓ | Get click count for a slug |

```bash
curl http://localhost:8080/api/stats/my-link \
  -H "Authorization: Bearer <your-key>"
```

```json
{ "slug": "my-link", "click_count": 42 }
```

---

### Health

```bash
curl http://localhost:8080/health
# ok
```

---

## Configuration

Configure via environment variables (copy `.env.example` to `.env`):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_URL` | `postgres://shortner:shortner@localhost:5432/shortner?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379/0` | Redis connection string |

---

## Project Structure

```
shortner-batch1/
├── cmd/server/          # Entry point
├── internal/
│   ├── model/           # Domain models (URL, APIKey, Stats)
│   ├── repository/      # PostgreSQL data access layer
│   ├── cache/           # Redis cache wrapper
│   ├── service/         # Business logic
│   ├── middleware/       # HTTP middleware (auth)
│   ├── handler/         # HTTP handlers
│   ├── server/          # Router assembly
│   └── static/          # Embedded web UI
├── migrations/          # SQL schema
├── web/                 # HTML source
├── docker-compose.yml
└── Dockerfile
```

---

## Development

### Run tests

```bash
go test ./...
```

All 43 unit tests run without any external dependencies (PostgreSQL/Redis are mocked).

### Build Docker image

```bash
docker build -t shortner:latest .
```

---

## Architecture

```
Browser / API Client
        │
        ▼
   HTTP Router (chi)
        │
   AuthMiddleware ──── APIKeyRepository ──── PostgreSQL
        │
   ┌────┴────────────────────┐
   │                         │
RedirectHandler          API Handlers
   │                    (URL / Stats / Auth)
   ├─ Cache.Get ──► Redis         │
   │                         Services
   └─ URLRepository ──► PostgreSQL    │
                                 Repositories ──► PostgreSQL
```

**Redirect flow** (optimized for speed):
1. Check Redis cache → hit: redirect immediately
2. Cache miss → query PostgreSQL → warm cache → redirect
3. Click count recorded asynchronously (goroutine, non-blocking)

---

## Roadmap

- [ ] Rate limiting on URL creation
- [ ] QR code generation per slug
- [ ] Detailed click analytics (referrer, country)
- [ ] API key scoping and rotation
- [ ] Admin dashboard
- [ ] `testcontainers-go` based integration tests

---

## Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Write tests first (TDD)
4. Commit your changes
5. Open a Pull Request

Please make sure all tests pass before submitting a PR:

```bash
go test ./... && go vet ./...
```

---

## License

MIT License — see [LICENSE](LICENSE) for details.
