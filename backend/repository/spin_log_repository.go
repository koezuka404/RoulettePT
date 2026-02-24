package repository

import (
	"backend/model"

	"gorm.io/gorm"
)

type ISpinLogRepository interface {
	Create(spinLog *model.SpinLog) error
	FindByUserIDAndIdempotencyKey(userID uint, idempotencyKey string) (*model.SpinLog, error)
}

type spinLogRepository struct {
	db *gorm.DB
}

func NewSpinLogRepository(db *gorm.DB) ISpinLogRepository {
	return &spinLogRepository{db}
}

func (r *spinLogRepository) Create(spinLog *model.SpinLog) error {
	return r.db.Create(spinLog).Error
}

func (r *spinLogRepository) FindByUserIDAndIdempotencyKey(userID uint, idempotencyKey string) (*model.SpinLog, error) {
	var log model.SpinLog
	if err := r.db.Where("user_id = ? AND idempotency_key = ?", userID, idempotencyKey).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}
