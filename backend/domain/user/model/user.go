package user

import "time"

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID           int64    `gorm:"primaryKey;autoIncrement"`
	Email        string   `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string   `gorm:"type:varchar(255);not null"`
	Role         UserRole `gorm:"type:varchar(10);not null"`

	TokenVersion int64 `gorm:"not null;default:0"`
	PointBalance int64 `gorm:"not null;default:0;check:point_balance >= 0"`
	IsActive     bool  `gorm:"not null;default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string { return "users" }
