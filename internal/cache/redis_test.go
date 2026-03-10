package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/bluejayA/shortner-batch1/internal/cache"
)

// Cache 인터페이스가 존재하는지 컴파일 타임 검증
var _ cache.Cache = (*mockCache)(nil)

type mockCache struct{}

func (m *mockCache) Get(_ context.Context, _ string) (string, error)                       { return "", nil }
func (m *mockCache) Set(_ context.Context, _, _ string, _ time.Duration) error             { return nil }
func (m *mockCache) Delete(_ context.Context, _ string) error                              { return nil }

func TestCacheMissError(t *testing.T) {
	// ErrCacheMiss가 정의되어 있는지 확인
	if cache.ErrCacheMiss == nil {
		t.Fatal("ErrCacheMiss가 nil임")
	}
}
