package main

import (
	"log"

	"roulettept/domain/models"
	"roulettept/infrastructure/db"
)

func main() {
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(
		&models.User{},
		&models.SpinLog{},
		&models.RefreshToken{},
		&models.PointAdjustment{},
		&models.AuditLog{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("migration success")
}
