package main

import (
	"context"
	"fmt"
	"go-auth-app/models"
	"go-auth-app/routes"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func accessSecretVersion(secretName string) string {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Printf("Failed to create secret manager client: %v\n", err)
		return ""
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		fmt.Printf("Failed to access secret: %v\n", err)
		return ""
	}

	return string(result.Payload.Data)
}

func loadSecrets() {
	os.Setenv("DB_HOST", accessSecretVersion("projects/706489728076/secrets/DB_HOST/versions/latest"))
	os.Setenv("DB_USER", accessSecretVersion("projects/706489728076/secrets/DB_USER/versions/latest"))
	os.Setenv("DB_NAME", accessSecretVersion("projects/706489728076/secrets/DB_NAME/versions/latest"))
	os.Setenv("DB_PASSWORD", accessSecretVersion("projects/706489728076/secrets/DB_PASSWORD/versions/latest"))
	os.Setenv("DB_PORT", accessSecretVersion("projects/706489728076/secrets/DB_PORT/versions/latest"))
	os.Setenv("DB_SSL", accessSecretVersion("projects/706489728076/secrets/DB_SSL/versions/latest"))
	os.Setenv("EMAIL_ADDRESS", accessSecretVersion("projects/706489728076/secrets/EMAIL_ADDRESS/versions/latest"))
	os.Setenv("EMAIL_PASSWORD", accessSecretVersion("projects/706489728076/secrets/EMAIL_PASSWORD/versions/latest"))
	os.Setenv("SMTP_HOST", accessSecretVersion("projects/706489728076/secrets/SMTP_HOST/versions/latest"))
	os.Setenv("SMTP_PORT", accessSecretVersion("projects/706489728076/secrets/SMTP_PORT/versions/latest"))
	os.Setenv("OPENAI_API_KEY", accessSecretVersion("projects/706489728076/secrets/OPENAI_API_KEY/versions/latest"))
}

func main() {
	loadSecrets()

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
