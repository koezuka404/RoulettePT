package userrepo

import (
	"context"

	user "roulettept/domain/user/model"
)

type UserRepository interface {
	// auth
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id int64) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	IncrementTokenVersion(ctx context.Context, userID int64) error

	// admin
	List(ctx context.Context, page, limit int, f UserListFilter) (items []user.User, total int64, err error)
	UpdateRole(ctx context.Context, userID int64, role user.UserRole) error
	Deactivate(ctx context.Context, userID int64) error

	// roulette（後で使う）
	AddPointsWithVersion(ctx context.Context, userID int64, expectedVersion int64, delta int64) (updated bool, err error)
}
