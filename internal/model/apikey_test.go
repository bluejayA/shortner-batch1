package model_test

import (
	"testing"

	"github.com/bluejayA/shortner-batch1/internal/model"
)

func TestAPIKey_HashKey(t *testing.T) {
	t.Run("동일한 키는 항상 동일한 해시를 반환한다", func(t *testing.T) {
		hash1 := model.HashKey("test-api-key-123")
		hash2 := model.HashKey("test-api-key-123")
		if hash1 != hash2 {
			t.Errorf("동일 키의 해시가 다름: %s != %s", hash1, hash2)
		}
	})

	t.Run("다른 키는 다른 해시를 반환한다", func(t *testing.T) {
		hash1 := model.HashKey("key-aaa")
		hash2 := model.HashKey("key-bbb")
		if hash1 == hash2 {
			t.Error("다른 키의 해시가 동일함")
		}
	})

	t.Run("해시 길이는 64자(SHA-256 hex)이다", func(t *testing.T) {
		hash := model.HashKey("any-key")
		if len(hash) != 64 {
			t.Errorf("해시 길이가 64가 아님: %d", len(hash))
		}
	})
}
