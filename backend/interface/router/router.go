package router

import "github.com/labstack/echo/v4"

func Register(e *echo.Echo, d Dependencies) {
	v1 := e.Group("/api/v1")

	RegisterAuth(v1, d)
	RegisterAdmin(v1, d)
	RegisterRoulette(v1, d)
}
