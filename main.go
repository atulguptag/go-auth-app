package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-auth-app/models"
	"go-auth-app/routes"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func accessSecretVersion(client *secretmanager.Client, secretName string) (string, error) {
	ctx := context.Background()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret %s: %v", secretName, err)
	}

	return string(result.Payload.Data), nil
}

func loadSecrets(client *secretmanager.Client) error {
	secrets := []struct {
		envName    string
		secretName string
	}{
		{"DB_HOST", "projects/706489728076/secrets/DB_HOST/versions/latest"},
		{"DB_USER", "projects/706489728076/secrets/DB_USER/versions/latest"},
		{"DB_NAME", "projects/706489728076/secrets/DB_NAME/versions/latest"},
		{"DB_PASSWORD", "projects/706489728076/secrets/DB_PASSWORD/versions/latest"},
		{"DB_PORT", "projects/706489728076/secrets/DB_PORT/versions/latest"},
		{"DB_SSL", "projects/706489728076/secrets/DB_SSL/versions/latest"},
		{"EMAIL_ADDRESS", "projects/706489728076/secrets/EMAIL_ADDRESS/versions/latest"},
		{"EMAIL_PASSWORD", "projects/706489728076/secrets/EMAIL_PASSWORD/versions/latest"},
		{"SMTP_HOST", "projects/706489728076/secrets/SMTP_HOST/versions/latest"},
		{"SMTP_PORT", "projects/706489728076/secrets/SMTP_PORT/versions/latest"},
		{"OPENAI_API_KEY", "projects/706489728076/secrets/OPENAI_API_KEY/versions/latest"},
	}

	for _, secret := range secrets {
		value, err := accessSecretVersion(client, secret.secretName)
		if err != nil {
			return err
		}
		os.Setenv(secret.envName, value)
	}

	return nil
}

func main() {
	isProduction := os.Getenv("GAE_ENV") == "standard"

	if isProduction {
		ctx := context.Background()
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create Secret Manager client: %v", err)
		}
		defer client.Close()

		if err := loadSecrets(client); err != nil {
			log.Fatalf("Error loading secrets: %v", err)
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, proceeding without it")
		} else {
			log.Println(".env file loaded successfully")
		}
	}

	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	config := models.Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
		SSLMode:  getEnvOrDefault("DB_SSL", "disable"),
	}

	if config.Password == "" {
		log.Fatal("Missing DB_PASSWORD environment variable")
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

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
