package model

import "time"

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID       uint     `json:"id" gorm:"primaryKey"`
	Email    string   `json:email" gorm:"unique;not null"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`

	TokenVersion int64 `json:"tokenVersion" gorm:"not null;default:0"`
	PointBalance int64 `json:"pointBalance" gorm:"not null;default:0;check:point_balance >= 0"`
	IsActive     bool  `json:"isActive" gorm:"not null;default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestBodyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ResponseUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"
	`
}

func (User) TableName() string { return "users" }
