package main

import (
	"log"

	"backend/db"
	user "backend/model"
)

func main() {
	// DB接続
	db := db.NewDB()

	// AutoMigrate
	if err := db.AutoMigrate(
		&user.User{},
		&user.RefreshToken{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	log.Println("migrate ok")
}
