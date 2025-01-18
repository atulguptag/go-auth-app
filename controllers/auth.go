package controllers

import (
	"fmt"
	"go-auth-app/models"
	"go-auth-app/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

// Login Function to authenticate a user
// c *gin.Context is a type that is used to get information about the incoming HTTP request and generate the HTTP response.
func Login(c *gin.Context) {
	var user models.User
	// ShouldBindJSON function binding the incoming JSON request body to a User struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// It queries the database for a user with the provided email.
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

	// Sets the expiration time of the token to 24 hours.
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &models.Claims{
		Email:  existingUser.Email,
		UserID: existingUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := access_token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	// SetCookie function sets a cookie in the response header.
	c.SetCookie("access_token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "Successfully logged in"})
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

	// Generate Verification Token
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &models.Claims{
		Email:  user.Email,
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating verification token"})
		return
	}

	verificationLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", tokenString)

	data := map[string]string{
		"VerificationLink": verificationLink,
	}

	go utils.SendEmail(user.Email, "Please Verify Your Email", "templates/email_verification_template.html", data)
	c.JSON(http.StatusOK, gin.H{"success": "User created successfully! Please check your email to verify your account."})
}

// Verify email function to verify an email address
func VerifyEmail(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
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
	c.JSON(200, gin.H{"success": "Email Verification Successful!"})
}

// Home Function to display home page
func Home(c *gin.Context) {
	//c.Cookie means it will get/read the cookie from the request.
	cookie, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(200, gin.H{"success": "Welcome to JokeMaster!", "email": claims.Email})
}

// Logout Function to logout a user
func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "Successfully logged out!"})
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
