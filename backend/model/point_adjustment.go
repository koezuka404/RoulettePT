package model

import "time"

type PointAdjustment struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index:idx_point_adj_user_created,priority:1"`
	AdminUserID uint      `gorm:"not null"`
	Delta       int64     `gorm:"not null"`
	Reason      string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"index:idx_point_adj_user_created,priority:2"`
}

func (PointAdjustment) TableName() string {
	return "point_adjustments"
}
