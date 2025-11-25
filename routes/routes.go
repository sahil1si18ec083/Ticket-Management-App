package routes

import (
	"ticket-app-gin-golang/controllers"
	"ticket-app-gin-golang/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	// Route registrations go here
	api := r.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUpController)
		auth.POST("/login", controllers.LoginController)

	}
	ticket := v1.Group("/ticket")
	ticket.Use(middleware.AuthMiddleware())
	{
		ticket.GET("/", controllers.GetUserTicketsController)
		ticket.GET("/:id", controllers.GetTicketByIDController)
		ticket.POST("/", controllers.CreateTicketController)
		ticket.DELETE("/:id", controllers.DeleteTicketByIDController)
		ticket.PUT("/:id", controllers.UpdateTicketByIDController)
	}
}
