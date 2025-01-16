package controllers

import (
	"go-auth-app/models"
	"go-auth-app/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("my_secret_key")

// Login Function to authenticate a user
// c *gin.Context is a type that is used to get information about the incoming HTTP request and generate the HTTP response.
func Login(c *gin.Context) {
	var user models.User
	// ShouldBindJSON function binding the incoming JSON request body to a User struct. If this fails, it returns a 400 error.
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// It queries the database for a user with the provided email. If no user found, it returns a 401 error.
	var existingUser models.User
	models.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID == 0 {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		c.JSON(400, gin.H{"error": "Invalid password!"})
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24) // Sets the expiration time of the token to 24 hours.

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := access_token.SignedString(jwtKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"})
		return
	}

	// SetCookie function sets a cookie in the response header. It takes the name, value, expiration time, path, domain, secure, and httponly as arguments.
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

	models.DB.Create(&user)
	c.JSON(200, gin.H{"success": "User created successfully"})
}

// Home Function to display home page
func Home(c *gin.Context) {
	cookie, err := c.Cookie("access_token") //c.Cookie means it will get/read the cookie from the request.
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, gin.H{"success": "home page", "role": claims.Role})
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
