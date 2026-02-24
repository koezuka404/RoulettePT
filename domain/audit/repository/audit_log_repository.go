package auditrepo

import (
	"context"

	audit "roulettept/domain/audit/model"
)

type AuditLogRepository interface {
	Create(ctx context.Context, l *audit.AuditLog) error
}
