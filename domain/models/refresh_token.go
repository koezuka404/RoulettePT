package models

import "time"

type RefreshToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement;column:id"`
	UserID    int64      `gorm:"not null;index:idx_refresh_tokens_user,priority:1;column:user_id"`
	TokenHash string     `gorm:"not null;uniqueIndex;column:token_hash"`
	ExpiresAt time.Time  `gorm:"not null;index:idx_refresh_tokens_expires;column:expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	UserAgent string     `gorm:"column:user_agent"`
	IP        string     `gorm:"column:ip"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime;column:created_at"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }
