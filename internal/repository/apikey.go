package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jay-ahn/shortner/internal/model"
)

// postgresAPIKeyRepository는 PostgreSQL 기반 APIKeyRepository 구현체다.
type postgresAPIKeyRepository struct {
	db *sql.DB
}

// NewAPIKeyRepository는 PostgreSQL APIKeyRepository를 생성한다.
func NewAPIKeyRepository(db *sql.DB) APIKeyRepository {
	return &postgresAPIKeyRepository{db: db}
}

func (r *postgresAPIKeyRepository) Insert(ctx context.Context, key *model.APIKey) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO api_keys (key_hash) VALUES ($1)`,
		key.KeyHash,
	)
	return err
}

func (r *postgresAPIKeyRepository) FindByHash(ctx context.Context, hash string) (*model.APIKey, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, key_hash, created_at, revoked_at FROM api_keys WHERE key_hash = $1`,
		hash,
	)
	key := &model.APIKey{}
	err := row.Scan(&key.ID, &key.KeyHash, &key.CreatedAt, &key.RevokedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("api key not found")
	}
	return key, err
}

func (r *postgresAPIKeyRepository) DeleteByHash(ctx context.Context, hash string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM api_keys WHERE key_hash = $1`,
		hash,
	)
	return err
}
