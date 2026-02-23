package gormrepo

import (
	"context"

	audit "roulettept/domain/audit/model"
	auditrepo "roulettept/domain/audit/repository"

	"gorm.io/gorm"
)

type auditLogRepo struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) auditrepo.AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, l *audit.AuditLog) error {
	return r.db.WithContext(ctx).Create(l).Error
}
