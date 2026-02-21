package repository

import (
	"context"
	"time"

	"roulettept/domain/models"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error)

	MarkUsed(ctx context.Context, id int64, usedAt time.Time) error
	DeleteByUserID(ctx context.Context, userID int64) error
	DeleteByHash(ctx context.Context, hash string) error

	DeleteExpired(ctx context.Context, now time.Time) (int64, error)
}
