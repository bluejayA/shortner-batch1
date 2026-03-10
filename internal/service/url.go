package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/jay-ahn/shortner/internal/cache"
	"github.com/jay-ahn/shortner/internal/model"
	"github.com/jay-ahn/shortner/internal/repository"
)

// ErrNotFound는 slug가 존재하지 않을 때 반환된다.
var ErrNotFound = errors.New("url not found")

// ErrExpired는 slug가 만료되었을 때 반환된다.
var ErrExpired = errors.New("url expired")

// ErrSlugConflict는 slug 충돌이 발생했을 때 반환된다.
var ErrSlugConflict = errors.New("slug already exists")

const (
	slugChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	slugLength = 6
	maxRetries = 3
	cacheTTL   = 24 * time.Hour
)

// urlService는 URLService의 구현체다.
type urlService struct {
	repo  repository.URLRepository
	cache cache.Cache
}

// NewURLService는 URLService를 생성한다.
func NewURLService(repo repository.URLRepository, c cache.Cache) URLService {
	return &urlService{repo: repo, cache: c}
}

// Create는 새 단축 URL을 생성한다. alias가 비어있으면 slug를 자동 생성한다.
func (s *urlService) Create(ctx context.Context, original, alias string, expiresAt *time.Time) (*model.URL, error) {
	slug := alias
	if slug == "" {
		var err error
		slug, err = s.generateUniqueSlug(ctx)
		if err != nil {
			return nil, err
		}
	}

	url := &model.URL{
		Slug:      slug,
		Original:  original,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Insert(ctx, url); err != nil {
		return nil, err
	}
	return url, nil
}

// Delete는 단축 URL을 삭제하고 캐시를 무효화한다.
func (s *urlService) Delete(ctx context.Context, slug string) error {
	if err := s.repo.DeleteBySlug(ctx, slug); err != nil {
		return err
	}
	// 캐시 무효화 (에러 무시 — 주요 작업 완료 후)
	_ = s.cache.Delete(ctx, slug)
	return nil
}

// Resolve는 slug에 해당하는 원본 URL을 반환한다. Redis 캐시를 우선 조회한다.
func (s *urlService) Resolve(ctx context.Context, slug string) (string, error) {
	// 캐시 우선 조회
	if original, err := s.cache.Get(ctx, slug); err == nil {
		return original, nil
	}

	// DB 조회
	url, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return "", ErrNotFound
	}
	if url.IsExpired() {
		return "", ErrExpired
	}

	// 캐시 저장 (만료일 기반 TTL)
	ttl := cacheTTL
	if url.ExpiresAt != nil {
		ttl = time.Until(*url.ExpiresAt)
	}
	_ = s.cache.Set(ctx, slug, url.Original, ttl)

	return url.Original, nil
}

// generateUniqueSlug는 충돌 없는 고유 slug를 생성한다.
func (s *urlService) generateUniqueSlug(ctx context.Context) (string, error) {
	for i := 0; i < maxRetries; i++ {
		slug := randomSlug()
		_, err := s.repo.FindBySlug(ctx, slug)
		if err != nil {
			// not found → 사용 가능
			return slug, nil
		}
	}
	return "", ErrSlugConflict
}

// randomSlug는 6자 alphanumeric 랜덤 문자열을 반환한다.
func randomSlug() string {
	b := make([]byte, slugLength)
	for i := range b {
		b[i] = slugChars[rand.Intn(len(slugChars))]
	}
	return string(b)
}
