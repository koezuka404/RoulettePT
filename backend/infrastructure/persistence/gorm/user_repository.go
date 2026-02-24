package gorm

import (
	"context"
	"errors"

	user "roulettept/domain/user/model"
	userrepo "roulettept/domain/user/repository"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userrepo.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) IncrementTokenVersion(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&user.User{}).
		Where("id = ?", userID).
		Update("token_version", gorm.Expr("token_version + 1")).Error
}

func (r *userRepo) AddPointsWithVersion(ctx context.Context, userID int64, expectedVersion int64, delta int64) (bool, error) {
	res := r.db.WithContext(ctx).
		Model(&user.User{}).
		Where("id = ? AND version = ?", userID, expectedVersion).
		Updates(map[string]any{
			"point_balance": gorm.Expr("point_balance + ?", delta),
			"version":       gorm.Expr("version + 1"),
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

// ------- admin -------

func (r *userRepo) List(ctx context.Context, page, limit int, f userrepo.UserListFilter) ([]user.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).Model(&user.User{})

	if f.Role != nil {
		q = q.Where("role = ?", string(*f.Role))
	}
	if f.IsActive != nil {
		q = q.Where("is_active = ?", *f.IsActive)
	}
	if f.Q != "" {
		q = q.Where("email ILIKE ?", f.Q+"%") // 仕様: 前方一致
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []user.User
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *userRepo) UpdateRole(ctx context.Context, userID int64, role user.UserRole) error {
	return r.db.WithContext(ctx).
		Model(&user.User{}).
		Where("id = ?", userID).
		Update("role", string(role)).Error
}

func (r *userRepo) Deactivate(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&user.User{}).
		Where("id = ?", userID).
		Update("is_active", false).Error
}
