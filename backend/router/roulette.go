package router

import (
	"backend/controller"

	"github.com/labstack/echo/v4"
)

// RegisterRoulette はルーレットルートを登録します（既存router.goは変更しない）
func RegisterRoulette(e *echo.Echo, rc controller.IRouletteController) {
	e.POST("/spin", rc.Spin)
}
