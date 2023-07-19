package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// UpdatePassword обрабатывает PUT-запросы на /profile/password
func UpdatePassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("userID").(uint)

	// Получите новый пароль из тела запроса
	var newPassword struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&newPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновите пароль пользователя в базе данных
	if err := db.Model(&User{}).Where("id = ?", userID).Update("password", newPassword.Password).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Password updated"})
}

func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

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

	router.Use(cors.New(config))

	// Add the db to the context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", RegisterUser)
		v1.POST("/login", LoginUser)
		v1.POST("/confirm", ConfirmRegistration)
		v1.POST("/resend_code", ResendConfirmationCode)
		v1.POST("/admin_password", SendAdminConfirmationCode)
		v1.POST("/admin_code", CheckAdminCode)
		v1.POST("/check-token-validity", CheckToken)
		v1.Use(AuthRequired())
		{
			v1.PUT("/profile/password", UpdatePassword)
		}
	}

	// Serve frontend static files
	router.StaticFS("/static", http.Dir("./frontend/build/static"))
	router.StaticFile("/", "./frontend/build/index.html")

	return router
}
