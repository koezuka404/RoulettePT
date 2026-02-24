package roulette

import (
	"context"
	"math/rand"
	"time"

	dmodel "roulettept/domain/roulette/model"
	drepo "roulettept/domain/roulette/repository"
)

type Service struct {
	logRepo  drepo.SpinLogRepository
	userRepo drepo.UserPointRepository
}

func New(
	logRepo drepo.SpinLogRepository,
	userRepo drepo.UserPointRepository,
) *Service {
	return &Service{
		logRepo:  logRepo,
		userRepo: userRepo,
	}
}

//////////////////////////////////////////////////////
// Spin
//////////////////////////////////////////////////////

func (s *Service) Spin(ctx context.Context, in SpinInput) (*SpinOutput, error) {

	if in.IdempotencyKey == "" {
		return nil, ErrInvalidKey
	}

	// idempotency確認
	existing, _ := s.logRepo.FindByKey(ctx, in.UserID, in.IdempotencyKey)
	if existing != nil {
		bal, err := s.userRepo.GetBalance(ctx, in.UserID)
		if err != nil {
			return nil, err
		}
		return &SpinOutput{
			Points:      existing.PointsEarned,
			NewBalance:  bal,
			IsDuplicate: true,
		}, nil
	}

	// 抽選
	points := drawReward()

	// 先にポイント加算（仕様必須順序）
	newBalance, err := s.userRepo.AddPoints(ctx, in.UserID, points)
	if err != nil {
		return nil, err
	}

	// 後からログ作成
	log := &dmodel.SpinLog{
		UserID:         in.UserID,
		IdempotencyKey: in.IdempotencyKey,
		PointsEarned:   points,
		CreatedAt:      time.Now(),
	}

	if err := s.logRepo.Create(ctx, log); err != nil {
		// 仕様通り：ログ失敗でもポイントは確定済み
		return nil, err
	}

	return &SpinOutput{
		Points:     points,
		NewBalance: newBalance,
	}, nil
}

//////////////////////////////////////////////////////
// History
//////////////////////////////////////////////////////

func (s *Service) History(ctx context.Context, in HistoryInput) (*HistoryOutput, error) {

	if in.Page <= 0 {
		return nil, ErrInvalidPage
	}

	if in.Limit <= 0 || in.Limit > 100 {
		return nil, ErrInvalidLimit
	}

	offset := (in.Page - 1) * in.Limit

	logs, total, err := s.logRepo.ListByUser(ctx, in.UserID, in.Limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]HistoryItem, 0, len(logs))

	for _, l := range logs {
		items = append(items, HistoryItem{
			ID:             l.ID,
			IdempotencyKey: l.IdempotencyKey,
			PointsEarned:   l.PointsEarned,
			CreatedAt:      l.CreatedAt.Format(time.RFC3339),
		})
	}

	return &HistoryOutput{
		Items: items,
		Total: total,
	}, nil
}

//////////////////////////////////////////////////////
// reward抽選
//////////////////////////////////////////////////////

func drawReward() int {
	r := rand.Intn(100)

	switch {
	case r < 50:
		return 10
	case r < 80:
		return 50
	case r < 95:
		return 100
	default:
		return 500
	}
}

func (s *Service) GetSpinHistory(ctx context.Context, in HistoryInput) (*HistoryOutput, error) {
	return s.History(ctx, in)
}
