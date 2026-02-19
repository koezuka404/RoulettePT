package dto

// --- Commands ---

type SpinCmd struct {
	UserID         uint
	IdempotencyKey string
}

type GetSpinHistoryCmd struct {
	UserID uint
	Page   int
	Limit  int
}

// --- Results ---

type SpinResult struct {
	SpinLogID     uint
	PointsEarned  int
	NewBalance    int
	RewardTier    string // normal/rare/super_rare
	IsDuplicate   bool
	SpunAtRFC3339 string
}

type SpinHistoryItem struct {
	SpinLogID     uint
	PointsEarned  int
	RewardTier    string
	SpunAtRFC3339 string
}

type SpinHistoryResult struct {
	Items []SpinHistoryItem
	Total int
	Page  int
	Limit int
}
