package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Logout で user_id を Context に入れる最小版
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := bearer(c.Request().Header.Get("Authorization"))
		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "missing token"))
		}

		secret := os.Getenv("SECRET")
		if secret == "" {
			return c.JSON(http.StatusInternalServerError, errResp("INTERNAL_SERVER_ERROR", "secret not set"))
		}

		claims := jwt.MapClaims{}
		t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// ✅ v5では Alg() で比較するのが安全
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !t.Valid {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "invalid token"))
		}

		sub, _ := claims["sub"].(string)
		uid, ok := parseInt64(sub)
		if !ok {
			return c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "invalid sub"))
		}

		c.Set("user_id", uid)
		return next(c)
	}
}

func bearer(hdr string) string {
	hdr = strings.TrimSpace(hdr)
	if strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
		return strings.TrimSpace(hdr[7:])
	}
	return ""
}

func parseInt64(s string) (int64, bool) {
	if s == "" {
		return 0, false
	}
	var n int64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		n = n*10 + int64(ch-'0')
	}
	return n, true
}
