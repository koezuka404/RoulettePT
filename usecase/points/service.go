package points

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	auditmodel "roulettept/domain/audit/model"
	auditrepo "roulettept/domain/audit/repository"
	pointsmodel "roulettept/domain/points/model"
	pointsrepo "roulettept/domain/points/repository"
	userrepo "roulettept/domain/user/repository"
)

var (
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrDeltaMustBeNonZero = errors.New("delta must be nonzero")
	ErrReasonRequired     = errors.New("reason required")

	// ビジネスルール違反（409想定）
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrBalanceOverflow     = errors.New("balance overflow")
	ErrConflict            = errors.New("conflict") // token_version 競合など
)

type Service struct {
	users userrepo.UserRepository
	adj   pointsrepo.PointAdjustmentRepository
	audit auditrepo.AuditLogRepository // nilでもOK
}

func NewService(
	users userrepo.UserRepository,
	adj pointsrepo.PointAdjustmentRepository,
	audit auditrepo.AuditLogRepository,
) *Service {
	return &Service{users: users, adj: adj, audit: audit}
}

func (s *Service) GetMyBalance(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, ErrInvalidUserID
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return u.PointBalance, nil
}

type AdminAdjustInput struct {
	AdminID int64
	UserID  int64
	Delta   int64
	Reason  string
}

func (s *Service) AdminAdjustPoints(ctx context.Context, in AdminAdjustInput) (newBalance int64, err error) {
	if in.AdminID <= 0 || in.UserID <= 0 {
		return 0, ErrInvalidUserID
	}
	if in.Delta == 0 {
		return 0, ErrDeltaMustBeNonZero
	}
	if strings.TrimSpace(in.Reason) == "" {
		return 0, ErrReasonRequired
	}

	// before（監査用 / expectedVersion 用）
	target, err := s.users.FindByID(ctx, in.UserID)
	if err != nil {
		return 0, err
	}
	beforeBalance := target.PointBalance
	expectedVersion := target.TokenVersion

	// 事前チェック（オーバーフロー / 0未満）
	if in.Delta > 0 && beforeBalance > math.MaxInt64-in.Delta {
		return 0, ErrBalanceOverflow
	}
	if in.Delta < 0 && beforeBalance+in.Delta < 0 {
		return 0, ErrInsufficientBalance
	}

	// ✅ ここが修正ポイント：AddPointsWithVersion を使う
	updated, err := s.users.AddPointsWithVersion(ctx, in.UserID, expectedVersion, in.Delta)
	if err != nil {
		return 0, err
	}
	if !updated {
		// token_version が更新されている等で競合した（同時更新）
		return 0, ErrConflict
	}

	// 更新後取得（シンプルに再取得）
	after, err := s.users.FindByID(ctx, in.UserID)
	if err != nil {
		return 0, err
	}
	newBalance = after.PointBalance

	// 調整履歴（失敗してもポイントは反映済みなので success 扱い）
	_ = s.adj.Create(ctx, &pointsmodel.PointAdjustment{
		UserID:      in.UserID,
		AdminUserID: in.AdminID,
		Delta:       in.Delta,
		Reason:      strings.TrimSpace(in.Reason),
	})

	// 監査ログ（任意）
	if s.audit != nil {
		beforeJSON, _ := json.Marshal(map[string]any{"point_balance": beforeBalance, "token_version": expectedVersion})
		afterJSON, _ := json.Marshal(map[string]any{"point_balance": newBalance, "token_version": after.TokenVersion})

		_ = s.audit.Create(ctx, &auditmodel.AuditLog{
			ActorUserID:  in.AdminID,
			Action:       "ADJUST_POINTS",
			ResourceType: "user",
			ResourceID:   in.UserID,
			BeforeJSON:   string(beforeJSON),
			AfterJSON:    string(afterJSON),
		})
	}

	return newBalance, nil
}
