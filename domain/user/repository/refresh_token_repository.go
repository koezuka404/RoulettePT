package userrepo

import (
	"context"
	"time"

	user "roulettept/domain/user/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *user.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*user.RefreshToken, error)
	MarkUsed(ctx context.Context, id int64, usedAt time.Time) error

	DeleteByHash(ctx context.Context, hash string) error
	DeleteByUserID(ctx context.Context, userID int64) error
}
