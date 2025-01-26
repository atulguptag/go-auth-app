package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-auth-app/models"
	"go-auth-app/utils"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Login Function to authenticate a user
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID == 0 {
		c.JSON(401, gin.H{"error": "User does not exists!"})
		return
	}

	if !existingUser.IsVerified {
		c.JSON(403, gin.H{"error": "Please verify your email address before logging in"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		c.JSON(400, gin.H{"error": "Invalid password!"})
		return
	}

	tokenString, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(200, gin.H{"success": "Successfully logged in", "access_token": tokenString})
}

// SignUp Function to create a new user
func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		c.JSON(409, gin.H{"error": "User already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		c.JSON(500, gin.H{"error": "Could not generate hash password"})
		return
	}

	user.IsVerified = false
	models.DB.Create(&user)

	// Generate JWT Token
	tokenString, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}

	verificationLink := fmt.Sprintf("https://jokemaster-go.netlify.app/verify?token=%s", tokenString)

	data := map[string]string{
		"VerificationLink": verificationLink,
	}
	templatePath := "templates/email_verification_template.html"
	go utils.SendEmail(user.Email, "Please Verify Your Email", templatePath, data)
	c.JSON(200, gin.H{"success": "User created successfully! Please check your email to verify your account."})
}

// Verify email function to verify an email address
func VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(400, gin.H{"error": "Token is required"})
		return
	}

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	models.DB.Where("email = ?", claims.Email).First(&user)
	if user.ID == 0 {
		c.JSON(400, gin.H{"error": "Invalid token"})
		return
	}

	user.IsVerified = true
	user.UpdatedAt = time.Now()
	models.DB.Save(&user)
	c.JSON(200, gin.H{"success": "Email Verification Successful! You can now login to your account."})
}

// Home Function to display home page
func Home(c *gin.Context) {
	// Read the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header missing"})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "Invalid authorization format"})
		return
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(200, gin.H{"success": "Welcome to JokeMaster!", "email": claims.Email})
}

// Logout Function to logout a user
func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": "Successfully logged out!",
	})
}

// ResetPassword Function to reset user password
func ResetPassword(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID == 0 {
		c.JSON(401, gin.H{"error": "User does not exist"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		c.JSON(500, gin.H{"error": "Could not generate hash password"})
		return
	}

	models.DB.Model(&existingUser).Update("password", user.Password)
	c.JSON(200, gin.H{"success": "Password reset successfully"})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var prompts []models.Prompt
	result := models.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&prompts)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve prompts"})
		return
	}

	c.JSON(200, prompts)
}

// GenerateState generates a random state string
func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GoogleLogin initiates the Google OAuth2 flow
func GoogleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate state"})
		return
	}

	url := utils.GetGoogleOAuthURL(state)
	c.Redirect(302, url)
}

// GoogleAuthCallback handles the callback from Google OAuth2
func GoogleAuthCallback(c *gin.Context) {
	stateFromQuery := c.Query("state")
	if stateFromQuery == "" {
		c.JSON(400, gin.H{"error": "State parameter missing"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "Code not found"})
		return
	}

	token, err := utils.ExchangeCode(code)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	client := utils.GoogleOauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user info"})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user info"})
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		user = models.User{
			Name:              userInfo.Name,
			Email:             userInfo.Email,
			IsVerified:        userInfo.VerifiedEmail,
			ImageURL:          userInfo.Picture,
			GoogleID:          userInfo.ID,
			Provider:          "google",
			VerificationToken: "",
			Prompts:           []models.Prompt{},
		}
		if err := models.DB.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate JWT Token
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate JWT"})
		return
	}

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, proceeding without it")
	} else {
		fmt.Println(".env file loaded successfully")
	}

	// Redirect to frontend with token
	frontendURL := os.Getenv("REACT_FRONTEND_URL")
	if frontendURL == "" {
		fmt.Println("REACT_FRONTEND_URL is not set in environment variables")
		c.JSON(500, gin.H{"error": "Configuration error"})
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/google/callback?token=%s", frontendURL, jwtToken)
	c.Redirect(302, redirectURL)
}
