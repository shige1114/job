package model

import "time"

// ApplicationDTO は Application 集約のデータを外部（インフラ層など）に渡すための単純な構造体です
type ApplicationDTO struct {
	ID        string
	Email     string
	Code      string
	ExpiresAt time.Time
}
