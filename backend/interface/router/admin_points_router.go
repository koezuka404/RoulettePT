package router

import (
	pointsctl "roulettept/interface/points/controller"
	"roulettept/usecase/points"

	"github.com/labstack/echo/v4"
)

func RegisterAdminPointsRoutes(g *echo.Group, svc *points.Service) {
	h := pointsctl.NewAdminPointsController(svc)
	g.POST("/points/adjust", h.Adjust)
}
