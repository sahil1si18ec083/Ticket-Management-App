package bootstrap

import (
	"fmt"
	"os"

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
	// Run versioned migrations (using gormigrate)
	err = Migrate(db)
	if err != nil {
		fmt.Println("Error during migrations:", err)
		os.Exit(1)
	}
	fmt.Println("Migrations completed successfully")
	return db

}
