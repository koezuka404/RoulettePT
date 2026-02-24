package dto

// POST /spin
type SpinRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
}

type SpinResponse struct {
	PointsEarned int64  `json:"points_earned"`
	NewBalance   int64  `json:"new_balance"`
	IsDuplicate  bool   `json:"is_duplicate"`
	RewardTier   string `json:"reward_tier,omitempty"`
}

// GET /spin/history
type SpinHistoryItem struct {
	ID             int64  `json:"id"`
	IdempotencyKey string `json:"idempotency_key"`
	PointsEarned   int64  `json:"points_earned"`
	CreatedAt      string `json:"created_at"` // RFC3339
}

type SpinHistoryResponse struct {
	Items []SpinHistoryItem `json:"items"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

// 共通レスポンス（必要なら）
type OKResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
