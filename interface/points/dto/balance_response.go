package dto

type BalanceResponse struct {
	Status string `json:"status"`
	Data   struct {
		PointBalance int64 `json:"point_balance"`
	} `json:"data"`
}
