package db

import (
	"em-test/cmd/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres provides a db-connection to Postgres using destination from caller
func ConnectPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot open db: %v", err)
	}

	if err := db.AutoMigrate(&model.Subscription{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	return db
}
