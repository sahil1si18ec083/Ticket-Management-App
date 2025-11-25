package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"ticket-app-gin-golang/models"
	"ticket-app-gin-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBInstance(db *gorm.DB) {

	DB = db
}

func SignUpController(c *gin.Context) {
	var reqbody models.SignupRequest
	var user models.User

	err := c.BindJSON(&reqbody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Request Body:", reqbody)
	res := DB.First(&user, "Email = ?", reqbody.Email)
	if res.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(reqbody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = hashedPassword
	user.Email = reqbody.Email
	user.Name = reqbody.Name
	tx := DB.Create(&user)
	fmt.Println("bbb", user.ID)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"token":   token,
	})

}

func LoginController(c *gin.Context) {
	var reqbody models.LoginRequest
	var user models.User

	err := c.BindJSON(&reqbody)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Request Body:", reqbody)
	res := DB.First(&user, "Email = ?", reqbody.Email)

	if res.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User credentials wrong"})
		return
	}

	// check whether password from client and password hashed in db are matching
	check := utils.CheckPasswordHash(reqbody.Password, user.Password)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User credentials wrong",
		})
		return
	}
	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating token"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": "login successfully",
		"token":   token,
	})
}
