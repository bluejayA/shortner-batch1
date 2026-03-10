package repository

import (
	"context"

	"github.com/bluejayA/shortner-batch1/internal/model"
)

// APIKeyRepository는 API 키 영속성 레이어 인터페이스다.
type APIKeyRepository interface {
	Insert(ctx context.Context, key *model.APIKey) error
	FindByHash(ctx context.Context, hash string) (*model.APIKey, error)
	DeleteByHash(ctx context.Context, hash string) error
}

// URLRepository는 URL 영속성 레이어 인터페이스다.
type URLRepository interface {
	Insert(ctx context.Context, url *model.URL) error
	FindBySlug(ctx context.Context, slug string) (*model.URL, error)
	DeleteBySlug(ctx context.Context, slug string) error
}

// StatsRepository는 클릭 통계 영속성 레이어 인터페이스다.
type StatsRepository interface {
	Increment(ctx context.Context, slug string) error
	FindBySlug(ctx context.Context, slug string) (*model.Stats, error)
}
