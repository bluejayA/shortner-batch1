package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bluejayA/shortner-batch1/internal/cache"
	"github.com/bluejayA/shortner-batch1/internal/model"
	"github.com/bluejayA/shortner-batch1/internal/service"
)

// mockURLRepo는 테스트용 URLRepository 목이다.
type mockURLRepo struct {
	insertFn      func(ctx context.Context, url *model.URL) error
	findBySlugFn  func(ctx context.Context, slug string) (*model.URL, error)
	deleteBySlugFn func(ctx context.Context, slug string) error
}

func (m *mockURLRepo) Insert(ctx context.Context, url *model.URL) error {
	return m.insertFn(ctx, url)
}
func (m *mockURLRepo) FindBySlug(ctx context.Context, slug string) (*model.URL, error) {
	return m.findBySlugFn(ctx, slug)
}
func (m *mockURLRepo) DeleteBySlug(ctx context.Context, slug string) error {
	return m.deleteBySlugFn(ctx, slug)
}

// mockCache는 테스트용 Cache 목이다.
type mockURLCache struct {
	getFn    func(ctx context.Context, slug string) (string, error)
	setFn    func(ctx context.Context, slug, url string, ttl time.Duration) error
	deleteFn func(ctx context.Context, slug string) error
}

func (m *mockURLCache) Get(ctx context.Context, slug string) (string, error) {
	return m.getFn(ctx, slug)
}
func (m *mockURLCache) Set(ctx context.Context, slug, url string, ttl time.Duration) error {
	return m.setFn(ctx, slug, url, ttl)
}
func (m *mockURLCache) Delete(ctx context.Context, slug string) error {
	return m.deleteFn(ctx, slug)
}

func TestURLService_Create(t *testing.T) {
	var stored *model.URL
	repo := &mockURLRepo{
		findBySlugFn: func(_ context.Context, _ string) (*model.URL, error) {
			return nil, errors.New("not found") // slug 충돌 없음
		},
		insertFn: func(_ context.Context, url *model.URL) error {
			stored = url
			return nil
		},
	}
	c := &mockURLCache{}
	svc := service.NewURLService(repo, c)

	url, err := svc.Create(context.Background(), "https://example.com", "", nil)
	if err != nil {
		t.Fatalf("Create() 에러: %v", err)
	}
	if url.Slug == "" {
		t.Fatal("slug가 비어있음")
	}
	if url.Original != "https://example.com" {
		t.Errorf("original이 다름: %s", url.Original)
	}
	if stored == nil {
		t.Fatal("Insert가 호출되지 않음")
	}
}

func TestURLService_Create_CustomAlias(t *testing.T) {
	repo := &mockURLRepo{
		findBySlugFn: func(_ context.Context, _ string) (*model.URL, error) {
			return nil, errors.New("not found")
		},
		insertFn: func(_ context.Context, _ *model.URL) error { return nil },
	}
	c := &mockURLCache{}
	svc := service.NewURLService(repo, c)

	url, err := svc.Create(context.Background(), "https://example.com", "my-alias", nil)
	if err != nil {
		t.Fatalf("Create() 에러: %v", err)
	}
	if url.Slug != "my-alias" {
		t.Errorf("커스텀 alias가 slug로 사용되지 않음: %s", url.Slug)
	}
}

func TestURLService_Delete(t *testing.T) {
	var deletedSlug, invalidatedSlug string
	repo := &mockURLRepo{
		deleteBySlugFn: func(_ context.Context, slug string) error {
			deletedSlug = slug
			return nil
		},
	}
	c := &mockURLCache{
		deleteFn: func(_ context.Context, slug string) error {
			invalidatedSlug = slug
			return nil
		},
	}
	svc := service.NewURLService(repo, c)

	if err := svc.Delete(context.Background(), "abc123"); err != nil {
		t.Fatalf("Delete() 에러: %v", err)
	}
	if deletedSlug != "abc123" {
		t.Errorf("삭제된 slug가 다름: %s", deletedSlug)
	}
	if invalidatedSlug != "abc123" {
		t.Errorf("캐시 무효화된 slug가 다름: %s", invalidatedSlug)
	}
}

func TestURLService_Resolve_CacheHit(t *testing.T) {
	dbCalled := false
	repo := &mockURLRepo{
		findBySlugFn: func(_ context.Context, _ string) (*model.URL, error) {
			dbCalled = true
			return nil, nil
		},
	}
	c := &mockURLCache{
		getFn: func(_ context.Context, _ string) (string, error) {
			return "https://cached.com", nil // 캐시 히트
		},
	}
	svc := service.NewURLService(repo, c)

	original, err := svc.Resolve(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("Resolve() 에러: %v", err)
	}
	if original != "https://cached.com" {
		t.Errorf("캐시 히트 값이 다름: %s", original)
	}
	if dbCalled {
		t.Error("캐시 히트인데 DB 조회됨")
	}
}

func TestURLService_Resolve_CacheMiss(t *testing.T) {
	var setCalled bool
	future := time.Now().Add(24 * time.Hour)
	repo := &mockURLRepo{
		findBySlugFn: func(_ context.Context, slug string) (*model.URL, error) {
			return &model.URL{Slug: slug, Original: "https://db.com", ExpiresAt: &future}, nil
		},
	}
	c := &mockURLCache{
		getFn: func(_ context.Context, _ string) (string, error) {
			return "", cache.ErrCacheMiss // 캐시 미스
		},
		setFn: func(_ context.Context, _, _ string, _ time.Duration) error {
			setCalled = true
			return nil
		},
	}
	svc := service.NewURLService(repo, c)

	original, err := svc.Resolve(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("Resolve() 에러: %v", err)
	}
	if original != "https://db.com" {
		t.Errorf("DB 조회 값이 다름: %s", original)
	}
	if !setCalled {
		t.Error("캐시 미스 후 캐시 저장이 호출되지 않음")
	}
}

func TestURLService_Resolve_Expired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	repo := &mockURLRepo{
		findBySlugFn: func(_ context.Context, slug string) (*model.URL, error) {
			return &model.URL{Slug: slug, Original: "https://old.com", ExpiresAt: &past}, nil
		},
	}
	c := &mockURLCache{
		getFn: func(_ context.Context, _ string) (string, error) {
			return "", cache.ErrCacheMiss
		},
	}
	svc := service.NewURLService(repo, c)

	_, err := svc.Resolve(context.Background(), "abc123")
	if !errors.Is(err, service.ErrExpired) {
		t.Errorf("만료 URL인데 ErrExpired가 아님: %v", err)
	}
}
