package repository

import (
	"context"
	"time"

	"roulettept/domain/models"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *models.RefreshToken) error
	FindByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	MarkUsed(ctx context.Context, id uint64, usedAt time.Time) error
	DeleteByUserID(ctx context.Context, userID uint64) error
	DeleteExpired(ctx context.Context, now time.Time) (int64, error)
}
