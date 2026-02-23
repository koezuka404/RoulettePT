package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"roulettept/infrastructure/db"
	gormrepo "roulettept/infrastructure/persistence/gorm"

	rouletteController "roulettept/interface/roulette/controller"
	rouletteRouter "roulettept/interface/router"
	"roulettept/usecase/roulette"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = godotenv.Load()

	gdb, err := db.NewDB()
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	// DI（Roulette）
	rouletteRepo := gormrepo.NewRouletteRepository(gdb)

	rouletteUC := roulette.New(rouletteRepo, rouletteRepo)
	rouletteHandler := rouletteController.New(rouletteUC)

	deps := rouletteRouter.Dependencies{
		Roulette: rouletteHandler,
	}

	// Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	e.GET("/readyz", func(c echo.Context) error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			return c.NoContent(http.StatusServiceUnavailable)
		}
		return c.NoContent(http.StatusOK)
	})

	api := e.Group("/api/v1")

	// Roulette router（router内で /roulette + JWT を付与）
	rouletteRouter.RegisterRoulette(api, deps)

	// Start / Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server start failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
