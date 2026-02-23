package middleware

import (
	"net/http"

	user "roulettept/domain/user/model"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleAny := c.Get("role")

		role, _ := roleAny.(user.UserRole)

		if role != user.RoleAdmin {
			return c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "admin only"))
		}

		return next(c)
	}
}
