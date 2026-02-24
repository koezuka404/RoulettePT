package model

import "time"

type SpinLog struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	UserID         int64     `gorm:"not null;uniqueIndex:uid_key"`
	IdempotencyKey string    `gorm:"type:varchar(128);not null;uniqueIndex:uid_key"`
	PointsEarned   int       `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
}

func (SpinLog) TableName() string {
	return "spin_logs"
}
