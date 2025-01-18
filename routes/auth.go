package routes

import (
	"go-auth-app/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/home", controllers.Home)
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.GET("/logout", controllers.Logout)
	r.POST("/reset-password", controllers.ResetPassword)
	r.POST("/generate-jokes", controllers.GenerateJokes)
	r.GET("/verify", controllers.VerifyEmail)
}
