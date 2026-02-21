package main

import (
	"log"

	"roulettept/domain/models"
	"roulettept/infrastructure/db"
)

func main() {
	// DB接続
	gdb, err := db.NewDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	// AutoMigrate（まずは User と RefreshToken だけ）
	if err := gdb.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	log.Println("migrate ok")
}
