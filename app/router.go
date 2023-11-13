
package app

// import the necessary packages/modules
import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// UpdatePassword is a function to handle PUT requests to /profile/password route
// This function updates the password of the user
func UpdatePassword(c *gin.Context) {
	// Get database instance and user ID from context
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// Retrieve the new password from the request body
	var newPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user's password in the database
	if err := db.Model(&User{}).Where("id = ?", userID).Update("password", newPassword.Password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Password updated"})
}

// SetupRoutes is a function to setup all the routes for the application
func SetupRoutes(db *gorm.DB) *gin.Engine {
	// Initialize router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	// Apply CORS to the router
	router.Use(cors.New(config))

	// Middleware to add the db to the context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", RegisterUser) // Route to register user
		v1.POST("/login", LoginUser) // Route to login user
		v1.POST("/confirm", ConfirmRegistration) // Route to confirm registration
		v1.POST("/resend_code", ResendConfirmationCode) // Route to resend confirmation code
		v1.POST("/admin_password", SendAdminConfirmationCode) // Route to send admin confirmation code
		v1.POST("/admin_code", CheckAdminCode) // Route to check admin code
		v1.POST("/check-token-validity", CheckToken) // Route to check token validity
		v1.Use(AuthRequired()) // Middleware to check auth token
		{
			v1.PUT("/profile/password", UpdatePassword) // Route to update password (requires auth)
		}
	}

	// Serve frontend static files
	router.StaticFS("/static", http.Dir("./frontend/build/static"))
	router.StaticFile("/", "./frontend/build/index.html")

	// Return the router
	return router
}