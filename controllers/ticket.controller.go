package controllers

import (
	"fmt"

	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	service *services.TicketService
}

func NewTicketController(service *services.TicketService) *TicketController {
	return &TicketController{
		service: service,
	}
}

// ----------------- Create -----------------
func (tc *TicketController) CreateTicket(c *gin.Context) {
	userID := c.GetUint("userID")

	var req models.TicketRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ticket, err := tc.service.CreateTicket(userID, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Ticket created successfully",
		"ticket":  ticket,
	})
}

// ----------------- Get All -----------------
func (tc *TicketController) GetUserTickets(c *gin.Context) {
	userID := c.GetUint("userID")

	role, exists := c.Get("role")
	fmt.Println(exists)
	fmt.Println("role is ", role)

	tickets, err := tc.service.GetUserTickets(userID, role.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	if len(tickets) == 0 {
		c.JSON(200, gin.H{"message": "No tickets found"})
		return
	}

	c.JSON(200, gin.H{
		"tickets": tickets,
	})
}

// ----------------- Get By ID -----------------
func (tc *TicketController) GetTicketByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	role, exists := c.Get("role")
	fmt.Println(exists)
	fmt.Println("role is ", role)

	ticket, err := tc.service.GetTicketByID(userID, id, role.(string))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ticket)
}

// ----------------- Update -----------------
func (tc *TicketController) UpdateTicket(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var req models.TicketUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ticket, err := tc.service.UpdateTicketByID(userID, id, req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Ticket updated successfully",
		"ticket":  ticket,
	})
}

// ----------------- Delete -----------------
func (tc *TicketController) DeleteTicket(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	role, exists := c.Get("role")
	fmt.Println(exists)
	fmt.Println("role is ", role)

	err := tc.service.DeleteTicketByID(userID, id, role.(string))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Ticket deleted successfully",
	})
}
