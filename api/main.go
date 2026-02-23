package main

import (
	"log"
	"os"

	"roulettept/infrastructure/db"
	gormrepo "roulettept/infrastructure/persistence/gorm"
	"roulettept/interface/router"
	usercontroller "roulettept/interface/user/controller"
	"roulettept/usecase/auth"
	"roulettept/usecase/useradmin"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	gdb, err := db.NewDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	userRepo := gormrepo.NewUserRepository(gdb)
	rtRepo := gormrepo.NewRefreshTokenRepository(gdb)
	auditRepo := gormrepo.NewAuditLogRepository(gdb)

	authSvc := auth.NewService(userRepo, rtRepo)
	adminSvc := useradmin.NewService(userRepo, rtRepo, auditRepo)

	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	authC := usercontroller.NewAuthController(authSvc)
	adminUserC := usercontroller.NewAdminUserController(adminSvc)

	router.Register(e, authC, adminUserC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}
