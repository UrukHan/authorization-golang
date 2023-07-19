package app

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

var AdminDataPath string

type AdminData struct {
	AdminPass       string `json:"admin_password"`
	AdminEmail      string `json:"admin_email"`
	TerminalKey     string `json:"terminal_key"`
	TerminalPass    string `json:"terminal_pass"`
	PolygonProvider string `json:"polygon_provider"`
	PolygonOwner    string `json:"polygon_owner"`
	PrivateKey      string `json:"private_key"`
}

func SendAdminConfirmationCode(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var json struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminData, err := ReadAdminDataFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read admin data"})
		return
	}

	if json.Password != adminData.AdminPass {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Неправильный пароль"})
		return
	}

	var admin User
	result := db.Where("email = ?", adminData.AdminEmail).First(&admin)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving admin data"})
		return
	}
	adminID := admin.ID

	adminToken, err := GenerateAdminToken(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating admin token"})
		return
	}

	emailBody := fmt.Sprintf("Ваш токен подтверждения: %s", adminToken)
	err = SendEmail(adminData.AdminEmail, "Токен подтверждения", emailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending confirmation token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Токен подтверждения отправлен на почту администратора"})
}

func CheckAdminCode(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing confirmation token"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation token, please try again."})
		return
	}

	if claims["role"] != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role in the token, please try again."})
		return
	}

	// Возвращаем полный токен для дальнейшей аутентификации
	c.JSON(http.StatusOK, gin.H{"token": input.Token})
}

func ReadAdminDataFromFile() (AdminData, error) {
	file, err := ioutil.ReadFile(AdminDataPath)

	if err != nil {
		return AdminData{}, err
	}

	data := AdminData{}

	err = json.Unmarshal([]byte(file), &data)

	if err != nil {
		return AdminData{}, err
	}

	return data, nil
}



