package router

import (
	"roulettept/interface/handler"
	appmw "roulettept/interface/middleware"

	"github.com/labstack/echo/v4"
)

type Deps struct {
	AuthHandler *handler.AuthHandler
	JWT         *appmw.JWTAuthMiddleware
}

func Register(e *echo.Echo, d Deps) {
	v1 := e.Group("/api/v1")

	// 仕様書は /healthz /readyz 分離推奨（まずは最小）
	v1.GET("/healthz", func(c echo.Context) error { return c.String(200, "ok") })
	v1.GET("/readyz", func(c echo.Context) error { return c.String(200, "ok") })

	auth := v1.Group("/auth")
	auth.POST("/register", d.AuthHandler.Register) // signup → register
	auth.POST("/login", d.AuthHandler.Login)

	// refresh/logout は次工程で CSRF + refresh_token cookie 実装後に追加
	// auth.POST("/refresh", d.AuthHandler.Refresh)
	// auth.POST("/logout", d.AuthHandler.Logout)

	_ = d // JWTは spin/points/admin を作る段階で使う
}
