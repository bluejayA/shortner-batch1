package service

import (
	"context"

	"github.com/bluejayA/shortner-batch1/internal/model"
	"github.com/bluejayA/shortner-batch1/internal/repository"
)

// statsService는 StatsService의 구현체다.
type statsService struct {
	repo repository.StatsRepository
}

// NewStatsService는 StatsService를 생성한다.
func NewStatsService(repo repository.StatsRepository) StatsService {
	return &statsService{repo: repo}
}

// Record는 slug의 클릭 수를 1 증가시킨다.
func (s *statsService) Record(ctx context.Context, slug string) error {
	return s.repo.Increment(ctx, slug)
}

// Get은 slug의 클릭 통계를 반환한다.
func (s *statsService) Get(ctx context.Context, slug string) (*model.Stats, error) {
	return s.repo.FindBySlug(ctx, slug)
}
