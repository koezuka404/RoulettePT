package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := ""
		if auth := c.Request().Header.Get("Authorization"); auth != "" {
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				tokenStr = strings.TrimSpace(auth[7:])
			}
		}
		if tokenStr == "" {
			if cookie, err := c.Cookie("token"); err == nil && cookie != nil && cookie.Value != "" {
				tokenStr = cookie.Value
			}
		}
		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		userIDVal, ok := claims["user_id"]
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		userIDFloat, ok := userIDVal.(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		userID := uint(userIDFloat)

		c.Set(UserIDKey, userID)
		return next(c)
	}
}
