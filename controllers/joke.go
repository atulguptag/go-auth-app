package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-auth-app/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JokeRequest struct {
	Prompt string `json:"prompt"`
}

type JokeResponse struct {
	English              []string `json:"english"`
	Hindi                []string `json:"hindi"`
	RemainingGenerations int      `json:"remaining_generations,omitempty"`
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateJokes(c *gin.Context) {
	var request JokeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	db, exists := c.Get("db")
	if !exists || db == nil {
		c.JSON(500, gin.H{"error": "Database not available"})
		return
	}
	dbConn := db.(*gorm.DB)

	// Determine if the request is from an authenticated user
	userID, authenticated := c.Get("userID")

	if !authenticated {
		handleAnonymousJokeGeneration(c, request, dbConn)
		return
	}

	handleAuthenticatedJokeGeneration(c, request, dbConn, userID.(uint))
}

func handleAnonymousJokeGeneration(c *gin.Context, request JokeRequest, db *gorm.DB) {
	anonymousID := c.GetHeader("X-Anonymous-Id")

	if anonymousID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing anonymous ID"})
		return
	}

	var anonymousGen models.AnonymousGeneration
	result := db.Where("anonymous_id = ?", anonymousID).First(&anonymousGen)

	// If no record exists or it's been more than 24 hours since last generation
	if result.Error == gorm.ErrRecordNotFound || time.Since(anonymousGen.LastGenerationTime) > 24*time.Hour {
		anonymousGen = models.AnonymousGeneration{
			AnonymousID:        anonymousID,
			GenerationCount:    1,
			LastGenerationTime: time.Now(),
		}
		if err := db.Create(&anonymousGen).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create anonymous generation record"})
			return
		}
	} else {
		// Check if generation limit is reached
		if anonymousGen.GenerationCount >= 3 {
			c.JSON(http.StatusForbidden, gin.H{
				"error":                 "You have reached the maximum number of free generations. Please sign up to continue.",
				"remaining_generations": 0,
			})
			return
		}

		anonymousGen.GenerationCount++
		anonymousGen.LastGenerationTime = time.Now()
		if err := db.Save(&anonymousGen).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update anonymous generation record"})
			return
		}
	}

	englishPrompt := fmt.Sprintf("Generate 5 funny jokes or puns based on these words: %s. Make them funny, creative, and humorous. Return only the jokes, one per line.", request.Prompt)
	englishJokes, err := callOpenAI(englishPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate English jokes"})
		return
	}

	hindiPrompt := fmt.Sprintf("Generate 5 funny jokes or puns in Hindi (using Devanagari script) based on these words: %s. Make them funny, creative, and humorous. Return only the jokes, one per line.", request.Prompt)
	hindiJokes, err := callOpenAI(hindiPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Hindi jokes"})
		return
	}

	response := JokeResponse{
		English:              parseJokes(englishJokes),
		Hindi:                parseJokes(hindiJokes),
		RemainingGenerations: 3 - anonymousGen.GenerationCount,
	}

	c.JSON(http.StatusOK, response)
}

func handleAuthenticatedJokeGeneration(c *gin.Context, request JokeRequest, db *gorm.DB, userID uint) {
	// Save the prompt to the database
	prompt := models.Prompt{
		UserID: userID,
		Text:   request.Prompt,
	}

	if err := db.Create(&prompt).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save the prompt"})
		return
	}

	englishPrompt := fmt.Sprintf("Generate 5 funny jokes or puns based on these words: %s. Make them funny, creative, and humorous. Return only the jokes, one per line.", request.Prompt)
	englishJokes, err := callOpenAI(englishPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate English jokes"})
		return
	}

	hindiPrompt := fmt.Sprintf("Generate 5 funny jokes or puns in Hindi (using Devanagari script) based on these words: %s. Make them funny, creative, and humorous. Return only the jokes, one per line.", request.Prompt)
	hindiJokes, err := callOpenAI(hindiPrompt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate Hindi jokes"})
		return
	}

	response := JokeResponse{
		English: parseJokes(englishJokes),
		Hindi:   parseJokes(hindiJokes),
	}

	c.JSON(200, response)
}

func callOpenAI(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")

	request := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func parseJokes(jokesString string) []string {
	// Split jokes by newline and filter empty lines
	var jokes []string
	for _, joke := range strings.Split(jokesString, "\n") {
		if trimmedJoke := strings.TrimSpace(joke); trimmedJoke != "" {
			jokes = append(jokes, trimmedJoke)
		}
	}
	return jokes
}
