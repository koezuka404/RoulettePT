package router

import (
	"roulettept/interface/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoulette(g *echo.Group, d Dependencies) {
	r := g.Group("/roulette", middleware.JWTMiddleware)

	r.POST("/spin", d.Roulette.Spin)
	r.GET("/history", d.Roulette.History)
}
