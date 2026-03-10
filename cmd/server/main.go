package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/bluejayA/shortner-batch1/internal/cache"
	"github.com/bluejayA/shortner-batch1/internal/handler"
	"github.com/bluejayA/shortner-batch1/internal/middleware"
	"github.com/bluejayA/shortner-batch1/internal/repository"
	"github.com/bluejayA/shortner-batch1/internal/server"
	"github.com/bluejayA/shortner-batch1/internal/service"
)

func main() {
	// 환경변수 읽기
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://shortner:shortner@localhost:5432/shortner?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/0")

	// PostgreSQL 연결
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("PostgreSQL 연결 실패: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("PostgreSQL ping 실패: %v", err)
	}

	// Redis 연결
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Redis URL 파싱 실패: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	// Repository 초기화
	urlRepo := repository.NewURLRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// Cache 초기화
	urlCache := cache.NewRedisCache(redisClient)

	// Service 초기화
	urlSvc := service.NewURLService(urlRepo, urlCache)
	authSvc := service.NewAuthService(apiKeyRepo)
	statsSvc := service.NewStatsService(statsRepo)

	// Handler 초기화
	redirectH := handler.NewRedirectHandler(urlSvc, statsSvc)
	urlH := handler.NewURLHandler(urlSvc)
	authH := handler.NewAuthHandler(authSvc)
	statsH := handler.NewStatsHandler(statsSvc)
	uiH := handler.NewUIHandler()
	authMw := middleware.NewAuthMiddleware(authSvc)

	// 라우터 조립
	r := server.NewRouter(redirectH, urlH, authH, statsH, authMw, uiH)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("서버 종료: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
