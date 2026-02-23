package router

import (
	"roulettept/interface/controller"
	"roulettept/interface/middleware"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, authC *controller.AuthController, adminC *controller.AdminUserController) {
	v1 := e.Group("/api/v1")

	// auth
	auth := v1.Group("/auth")
	auth.POST("/register", authC.Register)
	auth.POST("/login", authC.Login)
	auth.POST("/refresh", authC.Refresh, middleware.CSRFMiddleware)
	auth.POST("/logout", authC.Logout, middleware.JWTMiddleware, middleware.CSRFMiddleware)

	// admin users
	admin := v1.Group("/admin", middleware.JWTMiddleware, middleware.AdminOnly)
	admin.GET("/users", adminC.ListUsers)
	admin.PATCH("/users/:id/role", adminC.UpdateRole)
	admin.POST("/users/:id/deactivate", adminC.Deactivate)
	admin.POST("/users/:id/force-logout", adminC.ForceLogout)
}
