package gorm

import (
	"context"
	"errors"

	dmodel "roulettept/domain/roulette/model"
	drepo "roulettept/domain/roulette/repository"

	"gorm.io/gorm"
)

type rouletteRepository struct {
	db *gorm.DB
}

// コンパイル時に interface を満たしているか保証
var _ drepo.SpinLogRepository = (*rouletteRepository)(nil)
var _ drepo.UserPointRepository = (*rouletteRepository)(nil)

func NewRouletteRepository(db *gorm.DB) *rouletteRepository {
	return &rouletteRepository{db: db}
}

//////////////////////////////////////////////////////
// SpinLogRepository
//////////////////////////////////////////////////////

func (r *rouletteRepository) Create(ctx context.Context, log *dmodel.SpinLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *rouletteRepository) FindByKey(ctx context.Context, userID int64, key string) (*dmodel.SpinLog, error) {
	var log dmodel.SpinLog

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND idempotency_key = ?", userID, key).
		First(&log).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *rouletteRepository) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]dmodel.SpinLog, int64, error) {
	var logs []dmodel.SpinLog
	var total int64

	db := r.db.WithContext(ctx).
		Model(&dmodel.SpinLog{}).
		Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

//////////////////////////////////////////////////////
// UserPointRepository
//////////////////////////////////////////////////////

func (r *rouletteRepository) AddPoints(ctx context.Context, userID int64, delta int) (int64, error) {
	var user struct {
		PointBalance int64
	}

	err := r.db.WithContext(ctx).
		Table("users").
		Select("point_balance").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return 0, err
	}

	newBalance := user.PointBalance + int64(delta)
	if newBalance < 0 {
		return 0, errors.New("balance negative")
	}

	err = r.db.WithContext(ctx).
		Table("users").
		Where("id = ?", userID).
		Update("point_balance", newBalance).Error
	if err != nil {
		return 0, err
	}

	return newBalance, nil

}

func (r *rouletteRepository) GetBalance(ctx context.Context, userID int64) (int64, error) {
	var balance int64
	err := r.db.WithContext(ctx).
		Table("users").
		Select("point_balance").
		Where("id = ?", userID).
		Scan(&balance).Error
	if err != nil {
		return 0, err
	}
	return balance, nil
}
