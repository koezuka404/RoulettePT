package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// Refresh / Logout にだけ付ける前提（Double Submit + Origin/Referer）
func CSRFMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		csrfCookie, err := c.Cookie("csrf_token")
		if err != nil || csrfCookie.Value == "" {
			return c.JSON(http.StatusForbidden, errResp("CSRF_VALIDATION_FAILED", "csrf cookie missing"))
		}

		csrfHeader := strings.TrimSpace(c.Request().Header.Get("X-CSRF-Token"))
		if csrfHeader == "" || csrfHeader != csrfCookie.Value {
			return c.JSON(http.StatusForbidden, errResp("CSRF_VALIDATION_FAILED", "csrf mismatch"))
		}

		// Origin/Referer 簡易チェック（環境変数がある時だけ）
		allowed := os.Getenv("ALLOWED_ORIGINS") // 例: http://localhost:5173,https://xxx
		if allowed != "" {
			origin := c.Request().Header.Get("Origin")
			if origin == "" {
				origin = c.Request().Header.Get("Referer")
			}
			if origin != "" && !originAllowed(origin, allowed) {
				return c.JSON(http.StatusForbidden, errResp("CSRF_VALIDATION_FAILED", "origin not allowed"))
			}
		}

		return next(c)
	}
}

func originAllowed(origin, allowedCSV string) bool {
	for _, a := range strings.Split(allowedCSV, ",") {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if strings.HasPrefix(origin, a) {
			return true
		}
	}
	return false
}

type errResponse struct {
	Status string `json:"status"`
	Error  struct {
		Code    string        `json:"code"`
		Message string        `json:"message"`
		Details []interface{} `json:"details"`
	} `json:"error"`
}

func errResp(code, msg string) errResponse {
	var r errResponse
	r.Status = "error"
	r.Error.Code = code
	r.Error.Message = msg
	r.Error.Details = []interface{}{}
	return r
}
