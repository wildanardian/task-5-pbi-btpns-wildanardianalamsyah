package middlewares

import (
	"api-golang/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid or malformed token"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := helpers.VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or malformed token"})
			c.Abort()
			return
		}

		c.Set("user", claims.Email)
		c.Next()
	}

}
