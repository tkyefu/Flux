package database

import (
	"log"
	"flux/models"
)

// Migrate runs database migrations
func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Task{}, &models.PasswordReset{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}
