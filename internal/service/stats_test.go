package service_test

import (
	"context"
	"testing"

	"github.com/jay-ahn/shortner/internal/model"
	"github.com/jay-ahn/shortner/internal/service"
)

// mockStatsRepo는 테스트용 StatsRepository 목이다.
type mockStatsRepo struct {
	incrementFn  func(ctx context.Context, slug string) error
	findBySlugFn func(ctx context.Context, slug string) (*model.Stats, error)
}

func (m *mockStatsRepo) Increment(ctx context.Context, slug string) error {
	return m.incrementFn(ctx, slug)
}
func (m *mockStatsRepo) FindBySlug(ctx context.Context, slug string) (*model.Stats, error) {
	return m.findBySlugFn(ctx, slug)
}

func TestStatsService_Record(t *testing.T) {
	var recordedSlug string
	repo := &mockStatsRepo{
		incrementFn: func(_ context.Context, slug string) error {
			recordedSlug = slug
			return nil
		},
	}
	svc := service.NewStatsService(repo)

	if err := svc.Record(context.Background(), "abc123"); err != nil {
		t.Fatalf("Record() 에러: %v", err)
	}
	if recordedSlug != "abc123" {
		t.Errorf("기록된 slug가 다름: %s", recordedSlug)
	}
}

func TestStatsService_Get(t *testing.T) {
	repo := &mockStatsRepo{
		findBySlugFn: func(_ context.Context, slug string) (*model.Stats, error) {
			return &model.Stats{Slug: slug, ClickCount: 42}, nil
		},
	}
	svc := service.NewStatsService(repo)

	stats, err := svc.Get(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("Get() 에러: %v", err)
	}
	if stats.ClickCount != 42 {
		t.Errorf("클릭 수가 다름: %d", stats.ClickCount)
	}
	if stats.Slug != "abc123" {
		t.Errorf("slug가 다름: %s", stats.Slug)
	}
}
