package controllers

import (
	"github.com/ariefcatur/my-boilerplate/config"
	"github.com/ariefcatur/my-boilerplate/helpers"
	"github.com/ariefcatur/my-boilerplate/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Validation Error", gin.H{
			"details": err.Error(),
		})
		return
	}

	// Validate Email
	if !helpers.IsValidEmail(input.Email) {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Validation error", gin.H{
			"details": "Invalid email format",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Validation error", gin.H{
			"details": "Email already registered",
		})
	}

	// Check if username already exists
	if err := config.DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Validation error", gin.H{
			"details": "Username already taken",
		})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Server Error", gin.H{"details": "Failed to hash password"})
		return
	}

	// Create New User
	user := models.User{
		Username: input.Username,
		Email:    strings.ToLower(input.Email),
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Server Error", gin.H{"details": "Failed to create user"})
		return
	}

	helpers.APIResponse(c, http.StatusCreated, "User created successfully", gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Identity string `json:"identity" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Validation Error", gin.H{
			"details": err.Error(),
		})
		return
	}

	var user models.User

	// Check if identity is email or username
	query := config.DB
	if helpers.IsValidEmail(input.Identity) {
		query = query.Where("email = ?", strings.ToLower(input.Identity))
	} else {
		query = query.Where("username = ?", input.Identity)
	}

	if err := query.First(&user).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed", gin.H{
			"details": "Invalid credentials",
		})
		return
	}

	// Verificate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed", gin.H{
			"details": "Invalid credentials",
		})
		return
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Server error", gin.H{
			"details": "Failed to generate token",
		})
		return
	}

	helpers.APIResponse(c, http.StatusOK, "Login successful", gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
