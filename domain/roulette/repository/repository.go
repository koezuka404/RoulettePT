package repository

import (
	"context"

	model "roulettept/domain/roulette/model"
)

// SpinLogRepository: スピンログ永続化
type SpinLogRepository interface {
	Create(ctx context.Context, log *model.SpinLog) error
	FindByKey(ctx context.Context, userID int64, key string) (*model.SpinLog, error)
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]model.SpinLog, int64, error)
}

// UserPointRepository: ポイント加算（roulette用途）
type UserPointRepository interface {
	AddPoints(ctx context.Context, userID int64, delta int) (newBalance int64, err error)
}
