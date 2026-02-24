package dto

type AdminAdjustRequest struct {
	UserID int64  `json:"user_id"`
	Delta  int64  `json:"delta"`
	Reason string `json:"reason"`
}

type AdminAdjustResponse struct {
	Status string `json:"status"`
	Data   struct {
		NewBalance int64 `json:"new_balance"`
	} `json:"data"`
}
