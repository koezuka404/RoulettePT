package router

import "github.com/labstack/echo/v4"

func RegisterAuth(g *echo.Group, d Dependencies) {
	auth := g.Group("/auth")

	auth.POST("/register", d.Auth.Register)
	auth.POST("/login", d.Auth.Login)
	auth.POST("/refresh", d.Auth.Refresh)
	auth.POST("/logout", d.Auth.Logout)
}
