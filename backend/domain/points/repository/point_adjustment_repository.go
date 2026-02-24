package repository

import (
	"context"

	pointsmodel "roulettept/domain/points/model"
)

type PointAdjustmentRepository interface {
	Create(ctx context.Context, a *pointsmodel.PointAdjustment) error
}
