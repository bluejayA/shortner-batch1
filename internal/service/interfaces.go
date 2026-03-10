package service

import (
	"context"
	"time"

	"github.com/bluejayA/shortner-batch1/internal/model"
)

// AuthService는 API 키 인증 비즈니스 로직 인터페이스다.
type AuthService interface {
	Issue(ctx context.Context) (*model.APIKey, string, error)
	Revoke(ctx context.Context, key string) error
	Validate(ctx context.Context, key string) (bool, error)
}

// URLService는 URL 단축 비즈니스 로직 인터페이스다.
type URLService interface {
	Create(ctx context.Context, original, alias string, expiresAt *time.Time) (*model.URL, error)
	Delete(ctx context.Context, slug string) error
	Resolve(ctx context.Context, slug string) (string, error)
}

// StatsService는 클릭 통계 비즈니스 로직 인터페이스다.
type StatsService interface {
	Record(ctx context.Context, slug string) error
	Get(ctx context.Context, slug string) (*model.Stats, error)
}
