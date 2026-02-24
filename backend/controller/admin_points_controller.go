package controller

import (
	"net/http"

	usecase "backend/usecase"

	"github.com/labstack/echo/v4"
)

type IAdminPointsController interface {
	AdjustPoints(c echo.Context) error
}

type adminPointsController struct {
	pu usecase.IPointsUsecase
}

func NewAdminPointsController(pu usecase.IPointsUsecase) IAdminPointsController {
	return &adminPointsController{pu}
}

type AdminAdjustRequest struct {
	AdminID uint   `json:"admin_id"`
	UserID  uint   `json:"user_id"`
	Delta   int64  `json:"delta"`
	Reason  string `json:"reason"`
}

func (apc *adminPointsController) AdjustPoints(c echo.Context) error {
	var req AdminAdjustRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if req.AdminID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "admin_id is required"})
	}

	newBalance, err := apc.pu.AdminAdjustPoints(req.AdminID, req.UserID, req.Delta, req.Reason)
	if err != nil {
		switch err {
		case usecase.ErrDeltaZero:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "must_be_nonzero"})
		case usecase.ErrReasonRequired:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "required"})
		case usecase.ErrAdminRequired:
			return c.JSON(http.StatusForbidden, map[string]string{"error": "ADMIN_REQUIRED"})
		case usecase.ErrUserNotFoundPoints:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "NOT_FOUND"})
		case usecase.ErrInsufficientBalance:
			return c.JSON(http.StatusConflict, map[string]string{"error": "insufficient_balance"})
		case usecase.ErrBalanceOverflow:
			return c.JSON(http.StatusConflict, map[string]string{"error": "balance_overflow"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_balance": newBalance,
	})
}
