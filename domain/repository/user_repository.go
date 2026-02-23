package repository

import (
	"context"

	"roulettept/domain/models"
)

type UserListFilter struct {
	Role     *models.UserRole
	IsActive *bool
	Q        string
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id int64) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	IncrementTokenVersion(ctx context.Context, userID int64) error
	UpdateRole(ctx context.Context, userID int64, role models.UserRole) error
	Deactivate(ctx context.Context, userID int64) error

	List(ctx context.Context, page, limit int, f UserListFilter) ([]models.User, int64, error)
}
