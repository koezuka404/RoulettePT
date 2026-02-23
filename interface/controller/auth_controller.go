package controller

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"roulettept/usecase/auth"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	uc *auth.Service
}

func NewAuthController(uc *auth.Service) *AuthController {
	return &AuthController{uc: uc}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthController) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errResp("VALIDATION_ERROR", "invalid body"))
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errResp("VALIDATION_ERROR", "email and password required"))
	}

	if err := h.uc.SignUp(c.Request().Context(), auth.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		if err == auth.ErrEmailAlreadyUsed {
			return c.JSON(http.StatusConflict, errResp("CONFLICT", "already_taken"))
		}
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}

	return c.JSON(http.StatusOK, okResp(map[string]any{}))
}

func (h *AuthController) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errResp("VALIDATION_ERROR", "invalid body"))
	}

	out, err := h.uc.Login(c.Request().Context(), auth.LoginInput{
		Email:     strings.TrimSpace(req.Email),
		Password:  req.Password,
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	})
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "invalid credentials"))
		case auth.ErrUserInactive:
			return c.JSON(http.StatusForbidden, errResp("ACCOUNT_INACTIVE", "account inactive"))
		default:
			log.Printf("login error: %v", err)
			return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
		}
	}

	setRefreshCookie(c, out.RefreshToken)
	setCSRFCookie(c, out.CSRFToken)

	return c.JSON(http.StatusOK, okResp(map[string]any{
		"access_token": out.AccessToken,
	}))
}

func (h *AuthController) Refresh(c echo.Context) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil || refreshCookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, errResp("REFRESH_TOKEN_INVALID", "refresh token missing"))
	}

	access := bearer(c.Request().Header.Get("Authorization"))

	out, err := h.uc.Refresh(c.Request().Context(), auth.RefreshInput{
		RefreshToken: refreshCookie.Value,
		AccessToken:  access,
		UserAgent:    c.Request().UserAgent(),
		IP:           c.RealIP(),
	})
	if err != nil {
		switch err {
		case auth.ErrRefreshTokenInvalid:
			return c.JSON(http.StatusUnauthorized, errResp("REFRESH_TOKEN_INVALID", "refresh token invalid"))
		case auth.ErrRefreshTokenReused:
			return c.JSON(http.StatusUnauthorized, errResp("REFRESH_TOKEN_REUSED", "refresh token reused"))
		case auth.ErrTokenVersionMismatch:
			return c.JSON(http.StatusUnauthorized, errResp("TOKEN_VERSION_MISMATCH", "token version mismatch"))
		default:
			return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
		}
	}

	setRefreshCookie(c, out.RefreshToken)
	setCSRFCookie(c, out.CSRFToken)

	return c.JSON(http.StatusOK, okResp(map[string]any{
		"access_token": out.AccessToken,
	}))
}

func (h *AuthController) Logout(c echo.Context) error {
	refreshCookie, _ := c.Cookie("refresh_token")
	refresh := ""
	if refreshCookie != nil {
		refresh = refreshCookie.Value
	}

	uidAny := c.Get("user_id")
	userID, ok := uidAny.(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "unauthorized"))
	}

	if err := h.uc.Logout(c.Request().Context(), auth.LogoutInput{
		UserID:       userID,
		RefreshToken: refresh,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "server error"))
	}

	clearCookie(c, "refresh_token")
	clearCookie(c, "csrf_token")

	return c.JSON(http.StatusOK, okResp(map[string]any{}))
}

type okResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type errResponse struct {
	Status string `json:"status"`
	Error  struct {
		Code    string        `json:"code"`
		Message string        `json:"message"`
		Details []interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func okResp(data interface{}) okResponse {
	return okResponse{Status: "ok", Data: data}
}

func errResp(code, msg string) errResponse {
	var r errResponse
	r.Status = "error"
	r.Error.Code = code
	r.Error.Message = msg
	r.Error.Details = []interface{}{}
	return r
}

func bearer(hdr string) string {
	hdr = strings.TrimSpace(hdr)
	if strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
		return strings.TrimSpace(hdr[7:])
	}
	return ""
}

func setRefreshCookie(c echo.Context, refreshToken string) {
	secure := os.Getenv("COOKIE_SECURE") != "false"
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((14 * 24 * time.Hour).Seconds()),
	}
	c.SetCookie(cookie)
}

func setCSRFCookie(c echo.Context, csrf string) {
	secure := os.Getenv("COOKIE_SECURE") != "false"
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf,
		Path:     "/api/v1",
		HttpOnly: false,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((14 * 24 * time.Hour).Seconds()),
	}
	c.SetCookie(cookie)
}

func clearCookie(c echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/api/v1",
		MaxAge:   -1,
		HttpOnly: name == "refresh_token",
		Secure:   os.Getenv("COOKIE_SECURE") != "false",
		SameSite: http.SameSiteLaxMode,
	})
}
