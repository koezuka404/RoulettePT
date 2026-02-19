package models

import "time"

type AuditLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:id"`
	ActorUserID  int64     `gorm:"not null;index:idx_audit_logs_actor_created,priority:1;column:actor_user_id"`
	Action       string    `gorm:"not null;column:action"`
	ResourceType string    `gorm:"not null;column:resource_type"`
	ResourceID   int64     `gorm:"not null;column:resource_id"`
	BeforeJSON   string    `gorm:"column:before_json"`
	AfterJSON    string    `gorm:"column:after_json"`
	CreatedAt    time.Time `gorm:"not null;index:idx_audit_logs_actor_created,priority:2;autoCreateTime;column:created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
