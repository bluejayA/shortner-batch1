package model

import "time"

// Stats는 slug별 클릭 통계를 나타낸다.
type Stats struct {
	Slug       string
	ClickCount int64
	UpdatedAt  time.Time
}
