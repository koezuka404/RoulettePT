package controller

import (
	"net/http"
	"strconv"

	dto "roulettept/interface/roulette/dto"
	"roulettept/usecase/roulette"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	uc roulette.Usecase
}

func New(uc roulette.Usecase) *Handler {
	return &Handler{uc: uc}
}

// POST /spin
func (h *Handler) Spin(c echo.Context) error {
	var req dto.SpinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorPayload{Code: "INVALID_REQUEST", Message: "invalid body"},
		})
	}
	if req.IdempotencyKey == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorPayload{Code: "INVALID_REQUEST", Message: "idempotency_key required"},
		})
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorPayload{Code: "UNAUTHORIZED", Message: "missing user_id"},
		})
	}

	out, err := h.uc.Spin(c.Request().Context(), roulette.SpinInput{
		UserID:         userID,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OKResponse[dto.SpinResponse]{
		Status: "ok",
		Data: dto.SpinResponse{
			PointsEarned: int64(out.Points),
			NewBalance:   out.NewBalance,
			IsDuplicate:  out.IsDuplicate,
		},
	})
}

// GET /history
func (h *Handler) History(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorPayload{Code: "UNAUTHORIZED", Message: "missing user_id"},
		})
	}

	// デフォルト（usecase側が Page/Limit 必須なので）
	page := 1
	limit := 20

	if v := c.QueryParam("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			page = n
		}
	}
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	out, err := h.uc.GetSpinHistory(c.Request().Context(), roulette.HistoryInput{
		UserID: userID,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	items := make([]dto.SpinHistoryItem, 0, len(out.Items))
	for _, v := range out.Items {
		items = append(items, dto.SpinHistoryItem{
			PointsEarned: int64(v.Points),
			CreatedAt:    v.Time, // usecase側で RFC3339 string にしてる
		})
	}

	return c.JSON(http.StatusOK, dto.OKResponse[dto.SpinHistoryResponse]{
		Status: "ok",
		Data: dto.SpinHistoryResponse{
			Items: items,
			Total: out.Total,
			Page:  page,
			Limit: limit,
		},
	})
}
