package usecase

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"backend/model"
	"backend/repository"
)

var (
	ErrIdempotencyKeyRequired = errors.New("idempotency_key is required")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserInactive           = errors.New("account inactive")
	ErrPointBalanceOverflow   = errors.New("point balance overflow")
)

// 報酬テーブル: 10pt(50%), 50pt(30%), 100pt(15%), 500pt(5%)
var rewards = []struct {
	points int64
	weight int // 0-99の重み
}{
	{10, 50},
	{50, 30},
	{100, 15},
	{500, 5},
}

type IRouletteUsecase interface {
	Spin(userID uint, idempotencyKey string) (model.SpinLog, error)
}

type rouletteUsecase struct {
	users   repository.IUserRepositoryForRoulette
	spinLog repository.ISpinLogRepository
}

func NewRouletteUsecase(users repository.IUserRepositoryForRoulette, spinLog repository.ISpinLogRepository) IRouletteUsecase {
	return &rouletteUsecase{users: users, spinLog: spinLog}
}

func (u *rouletteUsecase) Spin(userID uint, idempotencyKey string) (model.SpinLog, error) {
	key := strings.TrimSpace(idempotencyKey)
	if key == "" {
		return model.SpinLog{}, ErrIdempotencyKeyRequired
	}

	// 冪等性: 同一キーで既に実行済みなら同じ結果を返す
	existing, err := u.spinLog.FindByUserIDAndIdempotencyKey(userID, key)
	if err != nil {
		return model.SpinLog{}, err
	}
	if existing != nil {
		return *existing, nil
	}

	// ユーザー取得・有効性チェック
	user, err := u.users.GetUserByID(userID)
	if err != nil {
		return model.SpinLog{}, err
	}
	if user == nil {
		return model.SpinLog{}, ErrUserNotFound
	}
	if !user.IsActive {
		return model.SpinLog{}, ErrUserInactive
	}

	// 1. ルーレット抽選
	points := drawReward()

	// 2. ポイント加算（排他制御）
	if err := u.users.AddPointsWithLock(userID, points); err != nil {
		return model.SpinLog{}, err
	}

	// 3. SpinLog作成
	spinLog := model.SpinLog{
		UserID:         userID,
		IdempotencyKey: key,
		PointsEarned:   points,
		CreatedAt:      time.Now(),
	}
	if err := u.spinLog.Create(&spinLog); err != nil {
		return model.SpinLog{}, err
	}

	return spinLog, nil
}

func drawReward() int64 {
	r := rand.Intn(100)
	acc := 0
	for _, rew := range rewards {
		acc += rew.weight
		if r < acc {
			return rew.points
		}
	}
	return rewards[0].points
}
