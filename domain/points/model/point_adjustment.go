package model

import "time"

// PointAdjustment ポイント手動調整履歴（管理者操作）
// 仕様: delta は負数OK。ただし残高が 0 未満にはならない / reason 必須 :contentReference[oaicite:2]{index=2}
type PointAdjustment struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	UserID      int64     `gorm:"not null;index"`
	AdminUserID int64     `gorm:"not null;index"`
	Delta       int64     `gorm:"not null"`
	Reason      string    `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

func (PointAdjustment) TableName() string { return "point_adjustments" }
