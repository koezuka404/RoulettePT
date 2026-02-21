package db

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

var _ repository.RefreshTokenRepository = (*refreshTokenRepo)(nil)

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

// used_at を NULL -> usedAt に更新（再利用/競合対策で条件付き）
func (r *refreshTokenRepo) MarkUsed(ctx context.Context, id int64, usedAt time.Time) error {
	res := r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", usedAt)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *refreshTokenRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).
		Error
}

func (r *refreshTokenRepo) DeleteByHash(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		Delete(&models.RefreshToken{}).
		Error
}

func (r *refreshTokenRepo) DeleteExpired(ctx context.Context, now time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", now).
		Delete(&models.RefreshToken{})

	return result.RowsAffected, result.Error
}
