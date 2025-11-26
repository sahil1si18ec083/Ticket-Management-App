package routes

import (
	"ticket-app-gin-golang/controllers"
	"ticket-app-gin-golang/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authController *controllers.AuthController, ticketController *controllers.TicketController) {

	api := r.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)
		auth.POST("/reset-password", authController.RequestPasswordReset)
	}

	ticket := v1.Group("/ticket")
	ticket.Use(middleware.AuthMiddleware())
	{
		ticket.POST("/", ticketController.CreateTicket)
		ticket.GET("/", ticketController.GetUserTickets)
		ticket.GET("/:id", ticketController.GetTicketByID)
		ticket.PUT("/:id", ticketController.UpdateTicket)
		ticket.DELETE("/:id", ticketController.DeleteTicket)
	}
}
