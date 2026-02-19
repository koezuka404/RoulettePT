package gorm

import (
	"context"
	"errors"

	"roulettept/domain/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateRole(ctx context.Context, id int64, role models.Role) (int64, error) {
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("role", role)
	return res.RowsAffected, res.Error
}

func (r *UserRepo) Deactivate(ctx context.Context, id int64) (int64, error) {
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false)
	return res.RowsAffected, res.Error
}

func (r *UserRepo) AddPoints(ctx context.Context, id int64, delta int) (int, int64, error) {
	// point_balance = point_balance + delta（原子的）
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("point_balance", gorm.Expr("point_balance + ?", delta))
	if res.Error != nil {
		return 0, 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, 0, nil
	}

	// 最新残高（最速優先で再取得）
	var u models.User
	if err := r.db.WithContext(ctx).Select("point_balance").First(&u, id).Error; err != nil {
		return 0, 0, err
	}
	return u.PointBalance, res.RowsAffected, nil
}

func (r *UserRepo) IncrementTokenVersion(ctx context.Context, id int64) (int, int64, error) {
	// token_version = token_version + 1（原子的）
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("token_version", gorm.Expr("token_version + 1"))
	if res.Error != nil {
		return 0, 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, 0, nil
	}

	// 最速優先で再取得
	var u models.User
	if err := r.db.WithContext(ctx).Select("token_version").First(&u, id).Error; err != nil {
		return 0, 0, err
	}
	return u.TokenVersion, res.RowsAffected, nil
}
