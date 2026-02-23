package repository

import (
	"context"

	"roulettept/domain/models"
)

type AuditLogRepository interface {
	Create(ctx context.Context, l *models.AuditLog) error
}
