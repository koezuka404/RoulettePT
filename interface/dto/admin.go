package dto

type ForceLogoutByAdminCmd struct {
	AdminID      uint
	TargetUserID uint
}

type ForceLogoutByAdminResult struct {
	TargetUserID             uint
	PreviousTokenVersion     int
	NewTokenVersion          int
	InvalidatedRefreshTokens int
	AuditLogID               uint
}

type UpdateUserRoleCmd struct {
	AdminID uint
	UserID  uint
	Role    string // USER/ADMIN
}

type UpdateUserRoleResult struct {
	UserID       uint
	PreviousRole string
	NewRole      string
	AuditLogID   uint
}

type DeactivateUserCmd struct {
	AdminID uint
	UserID  uint
}

type DeactivateUserResult struct {
	UserID                   uint
	IsActive                 bool
	InvalidatedRefreshTokens int
	AuditLogID               uint
}

type ListUsersCmd struct {
	Page     int
	Limit    int
	Role     *string
	IsActive *bool
	Q        *string // email prefix
}

type ListUsersItem struct {
	ID           uint
	Email        string
	Role         string
	PointBalance int
	IsActive     bool
}

type ListUsersResult struct {
	Items []ListUsersItem
	Total int
	Page  int
	Limit int
}
