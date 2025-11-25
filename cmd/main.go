package main

import (
	"fmt"
	"os"
	"ticket-app-gin-golang/controllers"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	} else {
		fmt.Println("Env file loaded successfully")
	}

}
func main() {
	initEnv()
	port := os.Getenv("PORT")
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Server starting on port:", port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to database successfully")

	// AutoMigrate User and Note models
	err = db.AutoMigrate(&models.User{}, &models.Ticket{})
	if err != nil {
		fmt.Println("Error during AutoMigrate:", err)
		os.Exit(1)
	}
	fmt.Println("AutoMigrate completed successfully")
	controllers.InitDBInstance(db)
	r := gin.Default()
	// p := ginprometheus.NewPrometheus("gin")
	// p.Use(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	routes.RegisterRoutes(r, db)

	r.Run(":" + port)
}
