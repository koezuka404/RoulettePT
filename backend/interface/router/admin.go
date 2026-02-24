package router

import (
	"roulettept/interface/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterAdmin(g *echo.Group, d Dependencies) {
	admin := g.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)

	admin.GET("/users", d.Admin.ListUsers)
	admin.PATCH("/users/:id/role", d.Admin.UpdateRole)
	admin.POST("/users/:id/deactivate", d.Admin.Deactivate)
	admin.POST("/users/:id/force-logout", d.Admin.ForceLogout)
}
