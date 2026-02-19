package models

import "time"

type SpinLog struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id"`
	UserID         int64     `gorm:"not null;index:idx_spin_logs_user_created,priority:1;uniqueIndex:uq_spin_logs_user_idem,priority:1;column:user_id"`
	IdempotencyKey string    `gorm:"not null;uniqueIndex:uq_spin_logs_user_idem,priority:2;column:idempotency_key"`
	PointsEarned   int       `gorm:"not null;column:points_earned"`
	CreatedAt      time.Time `gorm:"not null;index:idx_spin_logs_user_created,priority:2;autoCreateTime;column:created_at"`
}

func (SpinLog) TableName() string { return "spin_logs" }
