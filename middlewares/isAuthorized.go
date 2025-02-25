package middlewares

import (
	"go-auth-app/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsAuthorized(allowAnonymous bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			if allowAnonymous {
				c.Next()
				return
			}
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			if allowAnonymous {
				c.Next()
				return
			}

			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
			if allowAnonymous {
				c.Next()
				return
			}
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Parse the token
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
