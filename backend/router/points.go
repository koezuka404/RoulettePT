package router

import (
	"backend/controller"

	"github.com/labstack/echo/v4"
)

// RegisterPoints はポイントルートを登録します（既存router.goは変更しない）
func RegisterPoints(e *echo.Echo, pc controller.IPointsController, apc controller.IAdminPointsController) {
	e.GET("/points/balance", pc.GetBalance)
	e.POST("/admin/points/adjust", apc.AdjustPoints)
}
