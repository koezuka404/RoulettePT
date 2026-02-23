package middleware

import (
	"net/http"

	"roulettept/domain/models"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		roleAny := c.Get("role")
		role, _ := roleAny.(models.UserRole)
		if role != models.RoleAdmin {
			return c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "admin only"))
		}
		return next(c)
	}
}
