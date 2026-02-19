package dto

type GetMyBalanceCmd struct {
	UserID uint
}

type GetMyBalanceResult struct {
	UserID       uint
	PointBalance int
	AsOfRFC3339  string
}

type AdminAdjustPointsCmd struct {
	AdminID uint
	UserID  uint
	Delta   int
	Reason  string
}

type AdminAdjustPointsResult struct {
	UserID          uint
	PreviousBalance int
	Delta           int
	NewBalance      int
	Reason          string
	AdjustmentID    uint
	AuditLogID      uint
}
