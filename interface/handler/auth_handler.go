package handler

import (
	"net/http"

	"roulettept/interface/dto"
	"roulettept/interface/handler/response"
	"roulettept/usecase/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	uc auth.Usecase
}

func NewAuthHandler(uc auth.Usecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid body")
	}
	if req.Email == "" || req.Password == "" {
		return response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "email and password required",
			response.ErrorDetail{Field: "email/password", Issue: "required"},
		)
	}

	if err := h.uc.SignUp(c.Request().Context(), auth.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		// 最速優先：エラー分類は後で仕様書コードに寄せる
		return response.Fail(c, http.StatusBadRequest, "SIGNUP_FAILED", err.Error())
	}

	// 仕様書が data を求めるならここに入れる。今は最小で Created + ok
	return response.Created(c, map[string]any{"message": "registered"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid body")
	}
	if req.Email == "" || req.Password == "" {
		return response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "email and password required",
			response.ErrorDetail{Field: "email/password", Issue: "required"},
		)
	}

	out, err := h.uc.Login(c.Request().Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return response.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	return response.OK(c, dto.AuthResponse{AccessToken: out.AccessToken})
}
