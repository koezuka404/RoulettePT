package gorm

import (
	"context"
	"time"

	"roulettept/domain/models"
	"roulettept/domain/repository"

	"gorm.io/gorm"
)

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepo) FindByHash(ctx context.Context, hash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// ✅ interfaceに合わせて usedAt を受け取る
func (r *refreshTokenRepo) MarkUsed(ctx context.Context, id uint64, usedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", usedAt).
		Error
}

func (r *refreshTokenRepo) DeleteByUserID(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).
		Error
}

func (r *refreshTokenRepo) DeleteExpired(ctx context.Context, now time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&models.RefreshToken{})

	return result.RowsAffected, result.Error
}
