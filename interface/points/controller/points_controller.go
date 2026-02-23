package controller

import (
	"net/http"
	"strconv"

	"roulettept/interface/points/dto"
	"roulettept/usecase/points"

	"github.com/labstack/echo/v4"
)

type PointsController struct {
	uc *points.Service
}

func NewPointsController(uc *points.Service) *PointsController {
	return &PointsController{uc: uc}
}

func (h *PointsController) GetMyBalance(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok || userID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	bal, err := h.uc.GetMyBalance(c.Request().Context(), userID)
	if err != nil {
		// 不明エラーは 500（最短実装）
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var res dto.BalanceResponse
	res.Status = "ok"
	res.Data.PointBalance = bal
	return c.JSON(http.StatusOK, res)
}

// JWT middleware が c.Set("user_id", ...) している前提。
// 型が int / int64 / float64 / string のどれでも吸収する。
func getUserIDFromContext(c echo.Context) (int64, bool) {
	v := c.Get("user_id")
	switch t := v.(type) {
	case int64:
		return t, true
	case int:
		return int64(t), true
	case float64:
		// JSON経由などで稀にfloatになるケース救済
		return int64(t), true
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}
