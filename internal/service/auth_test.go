package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bluejayA/shortner-batch1/internal/model"
	"github.com/bluejayA/shortner-batch1/internal/service"
)

// mockAPIKeyRepo는 테스트용 APIKeyRepository 목(mock)이다.
type mockAPIKeyRepo struct {
	insertFn      func(ctx context.Context, key *model.APIKey) error
	findByHashFn  func(ctx context.Context, hash string) (*model.APIKey, error)
	deleteByHashFn func(ctx context.Context, hash string) error
}

func (m *mockAPIKeyRepo) Insert(ctx context.Context, key *model.APIKey) error {
	return m.insertFn(ctx, key)
}
func (m *mockAPIKeyRepo) FindByHash(ctx context.Context, hash string) (*model.APIKey, error) {
	return m.findByHashFn(ctx, hash)
}
func (m *mockAPIKeyRepo) DeleteByHash(ctx context.Context, hash string) error {
	return m.deleteByHashFn(ctx, hash)
}

func TestAuthService_Issue(t *testing.T) {
	var stored *model.APIKey
	repo := &mockAPIKeyRepo{
		insertFn: func(_ context.Context, key *model.APIKey) error {
			stored = key
			return nil
		},
	}
	svc := service.NewAuthService(repo)

	apiKey, plaintext, err := svc.Issue(context.Background())
	if err != nil {
		t.Fatalf("Issue() 에러: %v", err)
	}
	if plaintext == "" {
		t.Fatal("plaintext 키가 비어있음")
	}
	if apiKey.KeyHash == "" {
		t.Fatal("KeyHash가 비어있음")
	}
	if apiKey.KeyHash == plaintext {
		t.Fatal("KeyHash가 plaintext와 동일함 (해시 안됨)")
	}
	if stored == nil {
		t.Fatal("Insert가 호출되지 않음")
	}
	if stored.KeyHash != model.HashKey(plaintext) {
		t.Fatal("저장된 해시가 plaintext 해시와 다름")
	}
}

func TestAuthService_Revoke(t *testing.T) {
	t.Run("존재하는 키 폐기 성공", func(t *testing.T) {
		repo := &mockAPIKeyRepo{
			findByHashFn: func(_ context.Context, hash string) (*model.APIKey, error) {
				return &model.APIKey{KeyHash: hash}, nil
			},
			deleteByHashFn: func(_ context.Context, hash string) error {
				return nil
			},
		}
		svc := service.NewAuthService(repo)
		if err := svc.Revoke(context.Background(), "some-key"); err != nil {
			t.Fatalf("Revoke() 에러: %v", err)
		}
	})

	t.Run("존재하지 않는 키 폐기 시 에러", func(t *testing.T) {
		repo := &mockAPIKeyRepo{
			findByHashFn: func(_ context.Context, hash string) (*model.APIKey, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewAuthService(repo)
		if err := svc.Revoke(context.Background(), "no-such-key"); err == nil {
			t.Fatal("존재하지 않는 키 폐기 시 에러가 없음")
		}
	})
}

func TestAuthService_Validate(t *testing.T) {
	t.Run("유효한 키는 true 반환", func(t *testing.T) {
		repo := &mockAPIKeyRepo{
			findByHashFn: func(_ context.Context, hash string) (*model.APIKey, error) {
				return &model.APIKey{KeyHash: hash}, nil
			},
		}
		svc := service.NewAuthService(repo)
		ok, err := svc.Validate(context.Background(), "valid-key")
		if err != nil || !ok {
			t.Fatalf("유효한 키가 false 반환: ok=%v err=%v", ok, err)
		}
	})

	t.Run("폐기된 키는 false 반환", func(t *testing.T) {
		revokedAt := time.Now().Add(-1 * time.Hour)
		repo := &mockAPIKeyRepo{
			findByHashFn: func(_ context.Context, hash string) (*model.APIKey, error) {
				return &model.APIKey{KeyHash: hash, RevokedAt: &revokedAt}, nil
			},
		}
		svc := service.NewAuthService(repo)
		ok, err := svc.Validate(context.Background(), "revoked-key")
		if err != nil || ok {
			t.Fatalf("폐기된 키가 true 반환: ok=%v err=%v", ok, err)
		}
	})

	t.Run("존재하지 않는 키는 false 반환", func(t *testing.T) {
		repo := &mockAPIKeyRepo{
			findByHashFn: func(_ context.Context, hash string) (*model.APIKey, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewAuthService(repo)
		ok, err := svc.Validate(context.Background(), "no-key")
		if err != nil || ok {
			t.Fatalf("없는 키가 true 반환: ok=%v err=%v", ok, err)
		}
	})
}
