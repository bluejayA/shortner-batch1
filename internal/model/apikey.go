package model

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// APIKey는 API 인증 키 엔티티를 나타낸다.
type APIKey struct {
	ID        int
	KeyHash   string
	CreatedAt time.Time
	RevokedAt *time.Time
}

// IsRevoked는 키가 폐기되었는지 반환한다.
func (k *APIKey) IsRevoked() bool {
	return k.RevokedAt != nil
}

// HashKey는 주어진 키를 SHA-256으로 해시하여 hex 문자열로 반환한다.
func HashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", sum)
}
