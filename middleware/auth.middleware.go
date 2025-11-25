package middleware

import (
	"net/http"
	"strconv"
	"ticket-app-gin-golang/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		// Convert claims.Subject (string) → uint
		userIDInt, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token subject",
			})
			return
		}

		// store as uint
		c.Set("userID", uint(userIDInt))

		c.Next()
	}
}
