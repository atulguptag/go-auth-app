package main

import (
	"fmt"
	"go-auth-app/models"
	"go-auth-app/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	config := models.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	if config.Password == "" {
		fmt.Println("Missing DB_PASSWORD environment variable")
		return
	}

	models.InitDB(config)

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://jokemaster-go.netlify.app", "https://golang-deploy-448219.uc.r.appspot.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.Use(func(c *gin.Context) {
		c.Set("db", models.GetDB())
		c.Next()
	})

	routes.AuthRoutes(r)
	r.Run(":" + port)
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
