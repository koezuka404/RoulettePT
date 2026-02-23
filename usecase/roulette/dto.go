package roulette

type SpinInput struct {
	UserID         int64
	IdempotencyKey string
}

type SpinOutput struct {
	Points      int   `json:"points"`
	NewBalance  int64 `json:"new_balance"`
	IsDuplicate bool  `json:"is_duplicate"`
}

type HistoryInput struct {
	UserID int64
	Page   int
	Limit  int
}

type HistoryItem struct {
	ID             int64  `json:"id"`
	IdempotencyKey string `json:"idempotency_key"`
	PointsEarned   int    `json:"points_earned"`
	CreatedAt      string `json:"created_at"`
}

type HistoryOutput struct {
	Items []HistoryItem `json:"items"`
	Total int64         `json:"total"`
}
