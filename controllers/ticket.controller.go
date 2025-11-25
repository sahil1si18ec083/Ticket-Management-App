package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"ticket-app-gin-golang/models"

	"github.com/gin-gonic/gin"
)

// DB instance is declared in auth.controller.go

// Duplicate InitDBInstance function removed

func CreateTicketController(c *gin.Context) {
	fmt.Println("controller called")
	var request models.TicketRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var ticket models.Ticket
	ticket.Title = request.Title
	ticket.Content = request.Content
	ticket.Status = request.Status

	userid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fmt.Printf("%T,%v", userid, userid)
	useridStr, ok := userid.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID type"})
		return
	}

	idInt, err := strconv.Atoi(useridStr)
	if err != nil {
		fmt.Println("Conversion error:", err)
		return
	}

	useridUint := uint(idInt) // Convert int to uint

	ticket.UserID = useridUint

	fmt.Println(ticket.UserID)

	// Save ticket to DB
	if err := DB.Create(&ticket).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket created successfully",
		"ticket":  ticket,
	})
}
func GetUserTicketsController(c *gin.Context) {
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var tickets []models.Ticket
	var ticketResponse []models.TicketResponse
	res := DB.Where("user_id = ?", userId).Find(&tickets)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}
	for _, val := range tickets {
		var obj models.TicketResponse
		obj.ID = val.ID
		obj.UserID = val.UserID
		obj.Title = val.Title
		obj.Content = val.Content
		obj.Status = val.Status
		obj.CreatedAt = val.CreatedAt
		obj.UpdatedAt = val.UpdatedAt
		ticketResponse = append(ticketResponse, obj)

	}
	if len(ticketResponse) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tickets found"})
		return
	}
	c.JSON(200, gin.H{"tickets": ticketResponse})

}
func GetTicketByIDController(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var ticket models.Ticket
	res := DB.Where("user_id = ? AND id = ?", userId, id).First(&ticket)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	var ticketResp models.TicketResponse
	ticketResp.ID = ticket.ID
	ticketResp.UserID = ticket.UserID
	ticketResp.Title = ticket.Title
	ticketResp.Content = ticket.Content
	ticketResp.Status = ticket.Status
	ticketResp.CreatedAt = ticket.CreatedAt
	ticketResp.UpdatedAt = ticket.UpdatedAt
	c.JSON(http.StatusOK, ticketResp)

}
func UpdateTicketByIDController(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var ticket models.Ticket
	res := DB.Where("user_id = ? AND id = ?", userId, id).First(&ticket)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	var ticketUpdateRequest models.TicketUpdateRequest

	err := c.BindJSON(&ticketUpdateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(ticketUpdateRequest.Title)
	fmt.Println(ticketUpdateRequest.Content)
	if ticketUpdateRequest.Title != nil {
		ticket.Title = *ticketUpdateRequest.Title

	}
	if ticketUpdateRequest.Content != nil {
		ticket.Content = *ticketUpdateRequest.Content
	}
	if ticketUpdateRequest.Status != nil {
		ticket.Status = *ticketUpdateRequest.Status

	}
	result := DB.Save(&ticket)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return

	}
	var ticketResp models.TicketResponse
	ticketResp.ID = ticket.ID
	ticketResp.UserID = ticket.UserID
	ticketResp.Title = ticket.Title
	ticketResp.Content = ticket.Content
	ticketResp.Status = ticket.Status
	ticketResp.CreatedAt = ticket.CreatedAt
	ticketResp.UpdatedAt = ticket.UpdatedAt

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket updated successfully",
		"ticket":  ticketResp,
	})

}
func DeleteTicketByIDController(c *gin.Context) {
	id := c.Param("id")
	userId, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var ticket models.Ticket
	res := DB.Where("user_id = ? AND id = ?", userId, id).First(&ticket)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	if err := DB.Delete(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ticket"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})

}
