package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func GenerateConfirmationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(10000)
	return fmt.Sprintf("%04d", code)
}

func CheckToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header not provided"})
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header format must be Bearer <token>"})
		return
	}
	_, _, _, err := ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Invalid authorization token"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func GetUser(c *gin.Context, db *gorm.DB, userID uint) (User, error) {
	user := User{}
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return User{}, err
	}
	return user, nil
}

func GetNewOrder(db *gorm.DB) (uint, error) {
	var maxOrderID *uint

	// Выполняем запрос для получения максимального значения OrderID
	row := db.Table("transactions").Select("MAX(order_id)").Row() // укажите имя вашей таблицы вместо "transactions"
	err := row.Scan(&maxOrderID)

	if err != nil {
		return 250, errors.New("Ошибка при получении максимального значения OrderID")
	}

	if maxOrderID == nil {
		// нет строк в таблице или значений order_id, возвращаем 11
		fmt.Println("Нет строк в таблице или значений order_id")
		return 250, nil
	}

	// Увеличиваем значение OrderID на 1
	newOrderID := *maxOrderID + 1

	return newOrderID, nil
}


