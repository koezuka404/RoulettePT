package main

import (
	"log"

	audit "roulettept/domain/audit/model"
	user "roulettept/domain/user/model"
	"roulettept/infrastructure/db"
)

func main() {
	// DB接続
	gdb, err := db.NewDB()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	// AutoMigrate
	if err := gdb.AutoMigrate(
		&user.User{},
		&user.RefreshToken{},
		&audit.AuditLog{},
	); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	log.Println("migrate ok")
}
