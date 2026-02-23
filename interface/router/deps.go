package router

import (
	pointscontroller "roulettept/interface/points/controller"
	roulettecontroller "roulettept/interface/roulette/controller"
	usercontroller "roulettept/interface/user/controller"
)

type Dependencies struct {
	Auth     *usercontroller.AuthController
	Admin    *usercontroller.AdminUserController
	Roulette *roulettecontroller.Handler

	Points      *pointscontroller.PointsController
	AdminPoints *pointscontroller.AdminPointsController
}
