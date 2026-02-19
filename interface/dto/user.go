package dto

type UserView struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PointBalance int    `json:"point_balance"`
	IsActive     bool   `json:"is_active,omitempty"` // register 例にある
}
