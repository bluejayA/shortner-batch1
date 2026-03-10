package repository

import (
	"context"
	"database/sql"

	"github.com/bluejayA/shortner-batch1/internal/model"
)

// postgresStatsRepository는 PostgreSQL 기반 StatsRepository 구현체다.
type postgresStatsRepository struct {
	db *sql.DB
}

// NewStatsRepository는 PostgreSQL StatsRepository를 생성한다.
func NewStatsRepository(db *sql.DB) StatsRepository {
	return &postgresStatsRepository{db: db}
}

// Increment는 slug의 클릭 수를 1 증가시킨다. 행이 없으면 새로 삽입한다.
func (r *postgresStatsRepository) Increment(ctx context.Context, slug string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO click_stats (slug, click_count, updated_at)
		VALUES ($1, 1, NOW())
		ON CONFLICT (slug) DO UPDATE
		SET click_count = click_stats.click_count + 1,
		    updated_at  = NOW()
	`, slug)
	return err
}

func (r *postgresStatsRepository) FindBySlug(ctx context.Context, slug string) (*model.Stats, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT slug, click_count, updated_at FROM click_stats WHERE slug = $1`,
		slug,
	)
	s := &model.Stats{}
	if err := row.Scan(&s.Slug, &s.ClickCount, &s.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return &model.Stats{Slug: slug, ClickCount: 0}, nil
		}
		return nil, err
	}
	return s, nil
}
