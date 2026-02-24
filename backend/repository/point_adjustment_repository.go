package repository

import (
	"backend/model"

	"gorm.io/gorm"
)

type IPointAdjustmentRepository interface {
	Create(pa *model.PointAdjustment) error
	ListByUserID(userID uint) ([]model.PointAdjustment, error)
}

type pointAdjustmentRepository struct {
	db *gorm.DB
}

func NewPointAdjustmentRepository(db *gorm.DB) IPointAdjustmentRepository {
	return &pointAdjustmentRepository{db}
}

func (r *pointAdjustmentRepository) Create(pa *model.PointAdjustment) error {
	return r.db.Create(pa).Error
}

func (r *pointAdjustmentRepository) ListByUserID(userID uint) ([]model.PointAdjustment, error) {
	var list []model.PointAdjustment
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
