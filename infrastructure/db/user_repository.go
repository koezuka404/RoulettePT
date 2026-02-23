package db

import (
	"context"
	"strings"

	"roulettept/domain/models"
	"roulettept/domain/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

var _ repository.UserRepository = (*userRepository)(nil)

func NewUserRepository(database *gorm.DB) repository.UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) IncrementTokenVersion(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("token_version", gorm.Expr("token_version + 1")).
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

// List: ユーザー一覧 + 総件数
func (r *userRepository) List(ctx context.Context, page, limit int, f repository.UserListFilter) ([]models.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	qb := r.db.WithContext(ctx).Model(&models.User{})

	if f.Role != nil {
		qb = qb.Where("role = ?", *f.Role)
	}
	if f.IsActive != nil {
		qb = qb.Where("is_active = ?", *f.IsActive)
	}
	if s := strings.TrimSpace(f.Q); s != "" {
		like := "%" + s + "%"
		qb = qb.Where("email ILIKE ?", like)
	}

	var total int64
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	if err := qb.Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
