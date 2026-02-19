package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const CtxUserID = "user_id"

type JWTAuthMiddleware struct{}

func NewJWTAuthMiddleware() *JWTAuthMiddleware {
	return &JWTAuthMiddleware{}
}

func (m *JWTAuthMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secret := os.Getenv("SECRET")
		if secret == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "SECRET is not set")
		}

		auth := c.Request().Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing bearer token")
		}
		raw := strings.TrimPrefix(auth, "Bearer ")

		tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || tok == nil || !tok.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if ok {
			// user_id をcontextへ（型は環境差があるので float64→int64 対応）
			if v, ok := claims["user_id"]; ok {
				switch x := v.(type) {
				case float64:
					c.Set(CtxUserID, int64(x))
				case int64:
					c.Set(CtxUserID, x)
				}
			}
		}

		return next(c)
	}
}
