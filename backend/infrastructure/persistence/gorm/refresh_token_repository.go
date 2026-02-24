package gorm

import (
	"context"
	"errors"
	"time"

	user "roulettept/domain/user/model"
	userrepo "roulettept/domain/user/repository"

	"gorm.io/gorm"
)

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) userrepo.RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *user.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *refreshTokenRepo) FindByHash(ctx context.Context, hash string) (*user.RefreshToken, error) {
	var rt user.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepo) MarkUsed(ctx context.Context, id int64, usedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&user.RefreshToken{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", usedAt).Error
}

func (r *refreshTokenRepo) DeleteByHash(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		Delete(&user.RefreshToken{}).Error
}

func (r *refreshTokenRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&user.RefreshToken{}).Error
}
