package repository

import (
	"context"

	"roulettept/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id int64) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	UpdateRole(ctx context.Context, id int64, role models.Role) (rows int64, err error)
	Deactivate(ctx context.Context, id int64) (rows int64, err error)

	// Logout-all: token_version を +1（Txなし前提）
	IncrementTokenVersion(ctx context.Context, id int64) (newVersion int, rows int64, err error)

	AddPoints(ctx context.Context, id int64, delta int) (newBalance int, rows int64, err error)
}
