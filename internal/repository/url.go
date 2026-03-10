package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jay-ahn/shortner/internal/model"
)

// postgresURLRepository는 PostgreSQL 기반 URLRepository 구현체다.
type postgresURLRepository struct {
	db *sql.DB
}

// NewURLRepository는 PostgreSQL URLRepository를 생성한다.
func NewURLRepository(db *sql.DB) URLRepository {
	return &postgresURLRepository{db: db}
}

func (r *postgresURLRepository) Insert(ctx context.Context, url *model.URL) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO urls (slug, original, expires_at) VALUES ($1, $2, $3)`,
		url.Slug, url.Original, url.ExpiresAt,
	)
	return err
}

func (r *postgresURLRepository) FindBySlug(ctx context.Context, slug string) (*model.URL, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT slug, original, expires_at, created_at FROM urls WHERE slug = $1`,
		slug,
	)
	url := &model.URL{}
	err := row.Scan(&url.Slug, &url.Original, &url.ExpiresAt, &url.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("url not found")
	}
	return url, err
}

func (r *postgresURLRepository) DeleteBySlug(ctx context.Context, slug string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM urls WHERE slug = $1`, slug)
	return err
}
