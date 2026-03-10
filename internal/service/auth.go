package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/jay-ahn/shortner/internal/model"
	"github.com/jay-ahn/shortner/internal/repository"
)

// authService는 AuthService의 구현체다.
type authService struct {
	repo repository.APIKeyRepository
}

// NewAuthService는 AuthService를 생성한다.
func NewAuthService(repo repository.APIKeyRepository) AuthService {
	return &authService{repo: repo}
}

// Issue는 새 API 키를 발급하고 해시를 저장한다. plaintext는 1회만 반환된다.
func (s *authService) Issue(ctx context.Context) (*model.APIKey, string, error) {
	// 32바이트 랜덤 키 생성
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, "", err
	}
	plaintext := hex.EncodeToString(raw)

	key := &model.APIKey{
		KeyHash: model.HashKey(plaintext),
	}
	if err := s.repo.Insert(ctx, key); err != nil {
		return nil, "", err
	}
	return key, plaintext, nil
}

// Revoke는 주어진 키를 폐기(삭제)한다.
func (s *authService) Revoke(ctx context.Context, key string) error {
	hash := model.HashKey(key)
	if _, err := s.repo.FindByHash(ctx, hash); err != nil {
		return err
	}
	return s.repo.DeleteByHash(ctx, hash)
}

// Validate는 주어진 키가 유효한지 확인한다.
func (s *authService) Validate(ctx context.Context, key string) (bool, error) {
	hash := model.HashKey(key)
	apiKey, err := s.repo.FindByHash(ctx, hash)
	if err != nil {
		// not found 포함 — 유효하지 않음으로 처리
		return false, nil
	}
	return !apiKey.IsRevoked(), nil
}
