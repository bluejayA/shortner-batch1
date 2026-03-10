# ---- builder ----
FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shortner ./cmd/server

# ---- runner ----
FROM alpine:3.21 AS runner

RUN apk add --no-cache curl && \
    adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /build/shortner .

USER appuser

EXPOSE 8080

CMD ["./shortner"]
