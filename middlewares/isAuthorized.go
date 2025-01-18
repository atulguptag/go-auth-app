package middlewares

import (
	"go-auth-app/utils"

	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(cookie)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
