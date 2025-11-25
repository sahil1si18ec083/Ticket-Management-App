package bootstrap

import (
	"ticket-app-gin-golang/controllers"
	"ticket-app-gin-golang/repositories"
	"ticket-app-gin-golang/routes"
	"ticket-app-gin-golang/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitServer(db *gorm.DB) *gin.Engine {

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	ticketRepo := repositories.NewTicketRepository(db)
	ticketService := services.NewTicketService(ticketRepo)
	ticketController := controllers.NewTicketController(ticketService)
	r := gin.Default()
	routes.RegisterRoutes(r, authController, ticketController)
	return r

}
