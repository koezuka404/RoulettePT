package repository

import (
	"backend/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IUserRepositoryForRoulette はルーレット用のユーザーリポジトリインターフェース（既存IUserRepositoryは変更しない）
type IUserRepositoryForRoulette interface {
	GetUserByID(id uint) (*model.User, error)
	AddPointsWithLock(userID uint, points int64) error
}

// NewUserRepositoryForRoulette はルーレット用のユーザーリポジトリを返します
func NewUserRepositoryForRoulette(db *gorm.DB) IUserRepositoryForRoulette {
	return &userRepository{db: db}
}

// GetUserByID はIDでユーザーを1件取得します（ルーレット用）。
func (ur *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// AddPointsWithLock はユーザーのポイントを加算します（SELECT FOR UPDATEで排他制御）。
func (ur *userRepository) AddPointsWithLock(userID uint, points int64) error {
	return ur.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}
		return tx.Model(&user).Update("point_balance", gorm.Expr("point_balance + ?", points)).Error
	})
}
