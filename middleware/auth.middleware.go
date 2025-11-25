package middleware

import (
	"fmt"
	"net/http"

	"ticket-app-gin-golang/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("this is logging my auth middleware")
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token not provided",
			})
			return
		}
		claims, flag := utils.VerifyToken(token)
		if !flag {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token wrong",
			})
			return
		}
		c.Set("userID", claims.Subject)
		c.Next()

	}

}
