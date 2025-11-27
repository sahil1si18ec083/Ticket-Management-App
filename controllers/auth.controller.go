package controllers

import (
	"fmt"
	"net/http"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// -------------------
// Signup
// -------------------
func (ac *AuthController) Signup(c *gin.Context) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := ac.service.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"token":   token,
	})
}

// -------------------
// Login
// -------------------
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := ac.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func (ac *AuthController) ForgetPassword(c *gin.Context) {
	var req models.ForgetPasswordRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// check email is valid
	ac.service.ForgetPasswordReset(req.Email)

	c.JSON(200, gin.H{
		"message": "If the email exists, a reset link has been sent",
	})

}
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var req models.PasswordResetRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err = ac.service.ResetPassword(req.Email, req.TokenHash, req.NewPassword)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(200, gin.H{
		"message": "Password has been reset successfully",
	})

}
