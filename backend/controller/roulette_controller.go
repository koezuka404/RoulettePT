package controller

import (
	"net/http"

	usecase "backend/usecase"

	"github.com/labstack/echo/v4"
)

type IRouletteController interface {
	Spin(c echo.Context) error
}

type rouletteController struct {
	ru usecase.IRouletteUsecase
}

func NewRouletteController(ru usecase.IRouletteUsecase) IRouletteController {
	return &rouletteController{ru}
}

type SpinRequest struct {
	UserID         uint   `json:"user_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (rc *rouletteController) Spin(c echo.Context) error {
	var req SpinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if req.UserID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	spinLog, err := rc.ru.Spin(req.UserID, req.IdempotencyKey)
	if err != nil {
		switch err {
		case usecase.ErrIdempotencyKeyRequired:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "idempotency_key is required"})
		case usecase.ErrUserNotFound:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		case usecase.ErrUserInactive:
			return c.JSON(http.StatusForbidden, map[string]string{"error": "account inactive"})
		case usecase.ErrPointBalanceOverflow:
			return c.JSON(http.StatusConflict, map[string]string{"error": "point balance overflow"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"points_earned": spinLog.PointsEarned,
		"created_at":   spinLog.CreatedAt,
	})
}
