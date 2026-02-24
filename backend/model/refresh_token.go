package model

import "time"

type RefreshToken struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null"`

	TokenHash string `gorm:"not null;unique"`

	ExpiresAt time.Time `gorm:"not null"`

	UsedAt *time.Time

	UserAgent string
	IP        string

	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
