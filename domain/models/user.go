package models

import "time"

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Email        string    `gorm:"not null;uniqueIndex;column:email"`
	PasswordHash string    `gorm:"not null;column:password_hash"`
	Role         Role      `gorm:"not null;column:role"`
	TokenVersion int       `gorm:"not null;default:0;column:token_version"`
	PointBalance int       `gorm:"not null;default:0;check:point_balance >= 0;column:point_balance"`
	IsActive     bool      `gorm:"not null;default:true;column:is_active"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime;column:updated_at"`
}

func (User) TableName() string { return "users" }
