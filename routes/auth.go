package routes

import (
	"go-auth-app/controllers"
	"go-auth-app/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/home", controllers.Home)
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.GET("/logout", controllers.Logout)
	r.GET("/profile", middlewares.IsAuthorized(false), controllers.Profile)
	r.GET("/verify", controllers.VerifyEmail)
	r.POST("/reset-password", controllers.ResetPassword)
	r.POST("/generate-jokes", middlewares.IsAuthorized(true), controllers.GenerateJokes)
}
