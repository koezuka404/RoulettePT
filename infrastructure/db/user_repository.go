package db

import (
	"context"

	"roulettept/domain/models"
	drepo "roulettept/domain/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

var _ drepo.UserRepository = (*userRepository)(nil)

func NewUserRepository(database *gorm.DB) drepo.UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) IncrementTokenVersion(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("token_version", gorm.Expr("token_version + 1")).
		Error
}

func (r *userRepository) AddPoints(ctx context.Context, userID int64, delta int64) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("point_balance", gorm.Expr("point_balance + ?", delta)).
		Error
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int64, role models.UserRole) error {
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", role)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) Deactivate(ctx context.Context, userID int64) error {
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND is_active = true", userID).
		Update("is_active", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
