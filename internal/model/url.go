package model

import "time"

// URL은 단축 URL 엔티티를 나타낸다.
type URL struct {
	Slug      string
	Original  string
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// IsExpired는 만료일이 설정되어 있고 현재 시각이 만료일을 지났는지 반환한다.
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}
