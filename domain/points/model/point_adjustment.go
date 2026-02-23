package points

import "time"

type PointAdjustment struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;column:id"`
	UserID      int64     `gorm:"not null;index:idx_point_adjust_user_created,priority:1;column:user_id"`
	AdminUserID int64     `gorm:"not null;column:admin_user_id"`
	Delta       int       `gorm:"not null;column:delta"`
	Reason      string    `gorm:"not null;column:reason"`
	CreatedAt   time.Time `gorm:"not null;index:idx_point_adjust_user_created,priority:2;autoCreateTime;column:created_at"`
}

func (PointAdjustment) TableName() string { return "point_adjustments" }
