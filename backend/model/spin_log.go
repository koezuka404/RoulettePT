package model

import "time"

type SpinLog struct {
	ID             uint   `gorm:"primaryKey"`
	UserID         uint   `gorm:"not null;uniqueIndex:idx_spin_user_idempotency"`
	IdempotencyKey string `gorm:"not null;uniqueIndex:idx_spin_user_idempotency"`
	PointsEarned   int64  `gorm:"not null"`
	CreatedAt      time.Time
}

func (SpinLog) TableName() string {
	return "spin_logs"
}
