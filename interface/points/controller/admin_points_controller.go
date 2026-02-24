package controller

import (
	"net/http"
	"strings"

	"roulettept/interface/points/dto"
	"roulettept/usecase/points"

	"github.com/labstack/echo/v4"
)

type AdminPointsController struct {
	uc *points.Service
}

func NewAdminPointsController(uc *points.Service) *AdminPointsController {
	return &AdminPointsController{uc: uc}
}

func (h *AdminPointsController) Adjust(c echo.Context) error {
	adminID, ok := getUserIDFromContext(c)
	if !ok || adminID <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req dto.AdminAdjustRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	req.Reason = strings.TrimSpace(req.Reason)

	newBal, err := h.uc.AdminAdjustPoints(c.Request().Context(), points.AdminAdjustInput{
		AdminID: adminID,
		UserID:  req.UserID,
		Delta:   req.Delta,
		Reason:  req.Reason,
	})
	if err != nil {
		switch err {
		case points.ErrInvalidUserID, points.ErrDeltaMustBeNonZero, points.ErrReasonRequired:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case points.ErrInsufficientBalance, points.ErrBalanceOverflow, points.ErrConflict:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var res dto.AdminAdjustResponse
	res.Status = "ok"
	res.Data.NewBalance = newBal
	return c.JSON(http.StatusOK, res)
}
