package main

import (
	"fmt"
	"os"
	"ticket-app-gin-golang/controllers"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/repositories"
	"ticket-app-gin-golang/routes"
	"ticket-app-gin-golang/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	fmt.Println("Env file loaded successfully")
}

func main() {
	initEnv()

	port := os.Getenv("PORT")
	dsn := os.Getenv("DATABASE_URL")

	fmt.Println("Server starting on port:", port)

	// Connect DB
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

	// ----------------------------
	// DEPENDENCY INJECTION
	// ----------------------------

	// AUTH
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// TICKET
	ticketRepo := repositories.NewTicketRepository(db)
	ticketService := services.NewTicketService(ticketRepo)
	ticketController := controllers.NewTicketController(ticketService)

	// Init gin
	r := gin.Default()

	// Routes
	routes.RegisterRoutes(r, authController, ticketController)

	// Start server
	r.Run(":" + port)
}
