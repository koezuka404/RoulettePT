package main

import (
	"log"

	"roulettept/infrastructure/db"

	"github.com/joho/godotenv"
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

	// 例：Repo DI（使うなら）
	_ = db.NewUserRepository(gdb)
	_ = db.NewRefreshTokenRepository(gdb)

	// ルーター等はまだ無いなら、いったんここで終了でもOK（ビルド通すのが優先）
}
