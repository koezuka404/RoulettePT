package router

import (
	"net/http"
	"strings"

	"roulettept/interface/middleware"
	pointsctl "roulettept/interface/points/controller"

	"github.com/labstack/echo/v4"
)

func requireAdminRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("role").(string)
		if strings.ToUpper(role) != "ADMIN" {
			return echo.NewHTTPError(http.StatusForbidden, "admin only")
		}
		return next(c)
	}
}

func RegisterPointsRoutes(
	api *echo.Group, // ✅ *echo.Echo ではなく Group
	pointsController *pointsctl.PointsController,
	adminPointsController *pointsctl.AdminPointsController,
) {
	// user endpoints
	auth := api.Group("")
	auth.Use(middleware.JWTMiddleware)
	auth.GET("/points/balance", pointsController.GetMyBalance)

	// admin endpoints
	admin := api.Group("/admin")
	admin.Use(middleware.JWTMiddleware, requireAdminRole)
	admin.POST("/points/adjust", adminPointsController.Adjust)
}
