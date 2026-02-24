package controller

import (
	"net/http"
	"strconv"

	usecase "backend/usecase"

	"github.com/labstack/echo/v4"
)

type IPointsController interface {
	GetBalance(c echo.Context) error
}

type pointsController struct {
	pu usecase.IPointsUsecase
}

func NewPointsController(pu usecase.IPointsUsecase) IPointsController {
	return &pointsController{pu}
}

func (pc *pointsController) GetBalance(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	if userIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}
	userID64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}
	userID := uint(userID64)

	balance, err := pc.pu.GetMyBalance(userID)
	if err != nil {
		if err == usecase.ErrUserNotFoundPoints {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"point_balance": balance,
	})
}
