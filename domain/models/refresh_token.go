package models

import "time"

type RefreshToken struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	UserID int64 `gorm:"not null;index:idx_user_expires,priority:1"`

	TokenHash string `gorm:"not null;uniqueIndex"`

	ExpiresAt time.Time `gorm:"not null;index:idx_user_expires,priority:2"`

	UsedAt *time.Time

	UserAgent string
	IP        string

	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
