package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func authMilddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header request failed"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "authorizatio heasr format must be bearer <token>"})
			c.Abort()
			return
		}

		userID, err := validateToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}

}
