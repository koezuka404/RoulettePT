package gorm

import (
	"context"

	pointsmodel "roulettept/domain/points/model"
	pointsrepo "roulettept/domain/points/repository"

	"gorm.io/gorm"
)

type pointAdjustmentRepo struct {
	db *gorm.DB
}

func NewPointAdjustmentRepository(db *gorm.DB) pointsrepo.PointAdjustmentRepository {
	return &pointAdjustmentRepo{db: db}
}

func (r *pointAdjustmentRepo) Create(ctx context.Context, a *pointsmodel.PointAdjustment) error {
	return r.db.WithContext(ctx).Create(a).Error
}
