package app

import (
	"github.com/gin-gonic/gin"
	"net/http"http
	"strings"
)

// AuthRequired is a middleware function that checks for a valid authorization token.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header.
		authHeader := c.Request.Header.Get("Authorization")

		// If Authorization header is missing, return unauthorized status and stop the chain.
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}

		// Split the authHeader string into two parts.
		// The Authorization header should be in the format `Bearer <token>`.
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			return
		}

		// Parts[1] is expected to contain the singular token.
		// Parse the token parts to extract userID, adminID, and role.
		userID, adminID, role, err := ParseToken(parts[1])

		// If the token is invalid, return an unauthorized status.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}

		// If the token is valid, store the userID, adminID, and role in the context so they can be used in subsequent handlers.
		c.Set("userID", userID)
		c.Set("adminID", adminID)
		c.Set("role", role)

		// Continue processing the request.
		c.Next()
	}
}