package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss는 캐시에 키가 없을 때 반환되는 에러다.
var ErrCacheMiss = errors.New("cache miss")

// Cache는 slug → 원본 URL 캐싱 인터페이스다.
type Cache interface {
	Get(ctx context.Context, slug string) (string, error)
	Set(ctx context.Context, slug, url string, ttl time.Duration) error
	Delete(ctx context.Context, slug string) error
}

// redisCache는 Redis 기반 Cache 구현체다.
type redisCache struct {
	client *redis.Client
}

// NewRedisCache는 Redis Cache를 생성한다.
func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, slug string) (string, error) {
	val, err := c.client.Get(ctx, slug).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	}
	return val, err
}

func (c *redisCache) Set(ctx context.Context, slug, url string, ttl time.Duration) error {
	return c.client.Set(ctx, slug, url, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, slug string) error {
	return c.client.Del(ctx, slug).Err()
}
