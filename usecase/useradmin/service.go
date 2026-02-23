package useradmin

import (
	"context"
	"errors"
	"time"

	auditrepo "roulettept/domain/audit/repository"
	user "roulettept/domain/user/model"
	userrepo "roulettept/domain/user/repository"

	"gorm.io/gorm"
)

type Service struct {
	users userrepo.UserRepository
	rt    userrepo.RefreshTokenRepository
	audit auditrepo.AuditLogRepository
}

func NewService(
	users userrepo.UserRepository,
	rt userrepo.RefreshTokenRepository,
	audit auditrepo.AuditLogRepository,
) *Service {
	return &Service{
		users: users,
		rt:    rt,
		audit: audit, // 今は未使用でもOK（後で監査ログに使う）
	}
}

// controller が渡してくる入力
type ListInput struct {
	Page     int
	Limit    int
	Role     *user.UserRole
	IsActive *bool
	Q        string
}

// レスポンス（okResp(out) でそのまま返せる形）
type ListOutput struct {
	Users []UserSummary `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

type UserSummary struct {
	ID           int64         `json:"id"`
	Email        string        `json:"email"`
	Role         user.UserRole `json:"role"`
	TokenVersion int64         `json:"token_version"`
	PointBalance int64         `json:"point_balance"`
	IsActive     bool          `json:"is_active"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (s *Service) ListUsers(ctx context.Context, in ListInput) (ListOutput, error) {
	f := userrepo.UserListFilter{
		Role:     in.Role,
		IsActive: in.IsActive,
		Q:        in.Q,
	}

	usersList, total, err := s.users.List(ctx, in.Page, in.Limit, f)
	if err != nil {
		return ListOutput{}, err
	}

	out := make([]UserSummary, 0, len(usersList))
	for _, u := range usersList {
		out = append(out, UserSummary{
			ID:           u.ID,
			Email:        u.Email,
			Role:         u.Role,
			TokenVersion: u.TokenVersion,
			PointBalance: u.PointBalance,
			IsActive:     u.IsActive,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
		})
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	limit := in.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return ListOutput{
		Users: out,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *Service) UpdateRole(ctx context.Context, actorID, targetID int64, role user.UserRole) error {
	if targetID == 0 {
		return ErrNotFound
	}
	if actorID != 0 && actorID == targetID {
		return ErrSelfRoleChange
	}

	err := s.users.UpdateRole(ctx, targetID, role)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

func (s *Service) Deactivate(ctx context.Context, actorID, targetID int64) error {
	if targetID == 0 {
		return ErrNotFound
	}
	if actorID != 0 && actorID == targetID {
		return ErrSelfDeactivate
	}

	u, err := s.users.FindByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if u == nil {
		return ErrNotFound
	}
	if !u.IsActive {
		return ErrAlreadyInactive
	}

	if err := s.users.Deactivate(ctx, targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// 念のためトークンも失効
	_ = s.users.IncrementTokenVersion(ctx, targetID)
	_ = s.rt.DeleteByUserID(ctx, targetID)

	return nil
}

func (s *Service) ForceLogout(ctx context.Context, actorID, targetID int64) error {
	if targetID == 0 {
		return ErrNotFound
	}

	u, err := s.users.FindByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if u == nil {
		return ErrNotFound
	}

	// access token を全失効（token_version を上げる）
	if err := s.users.IncrementTokenVersion(ctx, targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// refresh token も全削除
	_ = s.rt.DeleteByUserID(ctx, targetID)

	return nil
}
