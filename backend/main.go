package main

import (
	"backend/controller"
	"backend/db"
	"backend/model"
	"backend/repository"
	"backend/router"
	"backend/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	//DBインスタンス化
	db := db.NewDB()
	if err := db.AutoMigrate(&model.SpinLog{}); err != nil {
		log.Fatalf("spin_log migrate failed: %v", err)
	}

	//userRepositoryインスタンス化
	userRepository := repository.NewUserRepository(db)
	//refreshReposityoryインスタンス化
	// refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	//userUsecaseインスタンス化
	userUsecase := usecase.NewUserUsecase(userRepository)
	//userControllerインスタンス化
	userController := controller.NewUserController(userUsecase)

	e := router.NewRouter(userController)
	spinLogRepository := repository.NewSpinLogRepository(db)
	userRepoForRoulette := repository.NewUserRepositoryForRoulette(db)
	rouletteUsecase := usecase.NewRouletteUsecase(userRepoForRoulette, spinLogRepository)
	rouletteController := controller.NewRouletteController(rouletteUsecase)
	router.RegisterRoulette(e, rouletteController)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	<-ctx.Done() // シグナル待ち
	e.Logger.Info("shutting down server...")
	e.Shutdown(context.Background())

}
