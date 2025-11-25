package bootstrap

import (
	"fmt"
	"os"
	"ticket-app-gin-golang/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database successfully")

	// AutoMigrate
	err = db.AutoMigrate(&models.User{}, &models.Ticket{})
	if err != nil {
		fmt.Println("Error during AutoMigrate:", err)
		os.Exit(1)
	}
	fmt.Println("AutoMigrate completed successfully")
	return db

}
