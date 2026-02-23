package db

import (
	"context"

	"roulettept/domain/models"
	"roulettept/domain/repository"

	"gorm.io/gorm"
)

type auditLogRepo struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, l *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(l).Error
}
