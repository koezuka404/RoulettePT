package usecase

import (
	"errors"
	"math"
	"strings"
	"time"

	"backend/model"
	"backend/repository"
)

var (
	ErrDeltaZero           = errors.New("must_be_nonzero")
	ErrReasonRequired      = errors.New("required")
	ErrAdminRequired       = errors.New("ADMIN_REQUIRED")
	ErrUserNotFoundPoints  = errors.New("NOT_FOUND")
	ErrInsufficientBalance = errors.New("insufficient_balance")
	ErrBalanceOverflow     = errors.New("balance_overflow")
)

type IPointsUsecase interface {
	GetMyBalance(userID uint) (int64, error)
	AdminAdjustPoints(adminID, userID uint, delta int64, reason string) (int64, error)
}

type pointsUsecase struct {
	users   repository.IUserRepositoryForRoulette
	adj     repository.IPointAdjustmentRepository
}

func NewPointsUsecase(users repository.IUserRepositoryForRoulette, adj repository.IPointAdjustmentRepository) IPointsUsecase {
	return &pointsUsecase{users: users, adj: adj}
}

func (u *pointsUsecase) GetMyBalance(userID uint) (int64, error) {
	user, err := u.users.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	if user == nil {
		return 0, ErrUserNotFoundPoints
	}
	return user.PointBalance, nil
}

func (u *pointsUsecase) AdminAdjustPoints(adminID, userID uint, delta int64, reason string) (int64, error) {
	if delta == 0 {
		return 0, ErrDeltaZero
	}
	if strings.TrimSpace(reason) == "" {
		return 0, ErrReasonRequired
	}

	admin, err := u.users.GetUserByID(adminID)
	if err != nil {
		return 0, err
	}
	if admin == nil || admin.Role != model.RoleAdmin {
		return 0, ErrAdminRequired
	}

	target, err := u.users.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	if target == nil {
		return 0, ErrUserNotFoundPoints
	}

	newBalance := target.PointBalance + delta
	if newBalance < 0 {
		return 0, ErrInsufficientBalance
	}
	if delta > 0 && target.PointBalance > math.MaxInt64-delta {
		return 0, ErrBalanceOverflow
	}

	if err := u.users.AddPointsWithLock(userID, delta); err != nil {
		return 0, err
	}

	pa := &model.PointAdjustment{
		UserID:      userID,
		AdminUserID: adminID,
		Delta:       delta,
		Reason:      strings.TrimSpace(reason),
		CreatedAt:   time.Now(),
	}
	if err := u.adj.Create(pa); err != nil {
		return 0, err
	}

	return newBalance, nil
}
