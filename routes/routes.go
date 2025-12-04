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
		auth.POST("/forget-password", authController.ForgetPassword)
		auth.POST("/reset-password", authController.ResetPassword)
	}

	tickets := v1.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.POST("/", ticketController.CreateTicket)
		tickets.GET("/", ticketController.GetUserTickets)
		tickets.GET("/:id", ticketController.GetTicketByID)
		tickets.PUT("/:id", ticketController.UpdateTicket)
		tickets.DELETE("/:id", ticketController.DeleteTicket)
		tickets.PUT("/:id/assign", ticketController.AssignTicket)
		tickets.PUT("/:id/unassign", ticketController.UnassignTicket)
	}
}
