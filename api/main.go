package main

import (
	"log"

	"roulettept/domain/models"
	"roulettept/domain/repository"
	"roulettept/infrastructure/db"
	persistencegorm "roulettept/infrastructure/persistence/gorm"
	"roulettept/interface/handler"
	appmw "roulettept/interface/middleware"
	"roulettept/interface/router"
	"roulettept/usecase/auth"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found (skip)")
	}

	gdb, err := db.NewDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	log.Println("DB connected")

	// 開発用 AutoMigrate（必要ならON）
	if err := gdb.AutoMigrate(
		&models.User{},
		&models.SpinLog{},
		&models.PointAdjustment{},
		&models.RefreshToken{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("AutoMigrate done")

	// --- DI ---
	var userRepo repository.UserRepository = persistencegorm.NewUserRepo(gdb)

	authSvc := auth.NewService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	jwtmw := appmw.NewJWTAuthMiddleware()

	// --- Echo ---
	e := echo.New()
	router.Register(e, router.Deps{
		AuthHandler: authHandler,
		JWT:         jwtmw,
	})

	log.Println("server starting on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
