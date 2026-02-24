package repository

import (
	"backend/model"

	"gorm.io/gorm"
)

type IRefreshTokenRepository interface {
	Create(refreshToken *model.RefreshToken) error
	// FindByHash(ctx context.Context, hash string) (*user.RefreshToken, error)
	// MarkUsed(ctx context.Context, id int64, usedAt time.Time) error

	// DeleteByHash(ctx context.Context, hash string) error
	// DeleteByUserID(ctx context.Context, userID int64) error
}
type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) IRefreshTokenRepository {
	return &refreshTokenRepository{db}
}

func (rtr *refreshTokenRepository) Create(refreshToken *model.RefreshToken) error {
	if err := rtr.db.Create(refreshToken).Error; err != nil {
		return err
	}
	return nil
}
