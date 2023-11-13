package app

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

// JwtSecret is a global variable to store the JWT Secret
var JwtSecret string

// GenerateToken is a function to create a JWT token for a user
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID, // Set user ID as a claim
		"exp":    time.Now().Add(time.Hour * 12).Unix(), // Set token expiration time
	}

	// Create new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret and return
	tokenString, err := token.SignedString([]byte(JwtSecret))
	return tokenString, err
}

// GenerateAdminToken is a function to create a JWT token for an admin
func GenerateAdminToken(adminID uint) (string, error) {
	claims := jwt.MapClaims{
		"adminID": adminID, // Set admin ID as a claim
		"role":    "admin", // Set role as a claim
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Set token expiration time
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret and return
	tokenString, err := token.SignedString([]byte(JwtSecret))
	return tokenString, err
}

// ParseToken is a function to validate and decode a JWT token
func ParseToken(tokenString string) (uint, uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JwtSecret), nil
	})

    // If there is an error, return it immediately
	if err != nil {
		return 0, 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
        // Check for token expiration
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return 0, 0, "", errors.New("Token expired")
		}

		var userID uint = 0
		var adminID uint = 0
		var role string = ""

        // Extract and return the claims from the token
		if claimUserID, ok := claims["userID"]; ok {
			userID = uint(claimUserID.(float64))
		}

		if claimAdminID, ok := claims["adminID"]; ok {
			adminID = uint(claimAdminID.(float64))
		}

		if claimRole, ok := claims["role"]; ok {
			role = claimRole.(string)
		}

		return userID, adminID, role, nil
	}
	return 0, 0, "", err
}