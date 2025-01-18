package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-auth-app/models"

	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JokeRequest struct {
	Prompt string `json:"prompt"`
}

type JokeResponse struct {
	English []string `json:"english"`
	Hindi   []string `json:"hindi"`
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

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Save the prompt to the database
	db := c.MustGet("db").(*gorm.DB)
	prompt := models.Prompt{
		UserID: userID.(uint),
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
