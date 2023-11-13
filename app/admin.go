package app

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// AdminPassword holds the admin password globally.
var AdminPassword string

// AdminEmail holds the admin's email globally.
var AdminEmail string

// SendAdminConfirmationCode is a function to validate the admin password, generate an admin token and send it to the admin email for identity confirmation. 
func SendAdminConfirmationCode(c *gin.Context) {
	// Get database instance from context
	db := c.MustGet("db").(*gorm.DB)
	
	// Declare JSON payload structure
	var json struct {
		Password string `json:"password" binding:"required"` // A required field in payload
		
	}

	// Bind incoming JSON to struct and check for errors
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password against global admin password
	if json.Password != AdminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Incorrect password"})
		return
	}

	// Retrieve the admin user from database with the global admin email
	var admin User
	result := db.Where("email = ?", AdminEmail).First(&admin)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving admin data"})
		return
	}
	adminID := admin.ID // Store admin ID

	// Generate admin token
	adminToken, err := GenerateAdminToken(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating admin token"})
		return
	}

	// Prepare and send confirmation email
	emailBody := fmt.Sprintf("Your confirmation token: %s", adminToken)
	err = SendEmail(AdminEmail, "Confirmation Token", emailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending confirmation token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation token sent to admin email"})
}

// CheckAdminCode is a function to validate the token received from an admin user. 
func CheckAdminCode(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"` // A required field in payload
	}

	// Bind incoming JSON to struct and check for errors
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing confirmation token"})
		return
	}

	// Parse and validate the token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})

	// Check for errors and token validity
	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation token, please try again."})
		return
	}

	// Ensure token belongs to admin
	if claims["role"] != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role in the token, please try again."})
		return
	}

	// Return full token for future authentication
	c.JSON(http.StatusOK, gin.H{"token": input.Token})
}