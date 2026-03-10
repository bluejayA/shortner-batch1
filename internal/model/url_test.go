package model_test

import (
	"testing"
	"time"

	"github.com/bluejayA/shortner-batch1/internal/model"
)

func TestURL_IsExpired(t *testing.T) {
	t.Run("만료일이 지난 URL은 true를 반환한다", func(t *testing.T) {
		past := time.Now().Add(-1 * time.Hour)
		u := model.URL{ExpiresAt: &past}
		if !u.IsExpired() {
			t.Error("만료된 URL이 IsExpired() = false를 반환함")
		}
	})

	t.Run("만료일이 미래인 URL은 false를 반환한다", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		u := model.URL{ExpiresAt: &future}
		if u.IsExpired() {
			t.Error("유효한 URL이 IsExpired() = true를 반환함")
		}
	})

	t.Run("만료일이 없는 URL은 false를 반환한다", func(t *testing.T) {
		u := model.URL{}
		if u.IsExpired() {
			t.Error("만료일 없는 URL이 IsExpired() = true를 반환함")
		}
	})
}
