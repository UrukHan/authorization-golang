package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = 587
	smtpUser   = "neurobet.mail@gmail.com"
	smtpPass   = "jomgzuajguxutrlb"
)

var MaxConfirmationTries int

type User struct {
	gorm.Model
	Phone             string
	Email             string `gorm:"unique"`
	Address           string
	Password          string
	AccessTo          time.Time
	Subscribe         string
	PaymentId         string
	Confirmed         bool
	ConfirmedCode     string
	ConfirmationTries int
	Block             bool
}

type UserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SetPassword hashes the password and stores it in the Password field
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

// CheckPassword compares the provided password against the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func RegisterUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
		return
	}

	var user User
	err := db.Where("email = ?", input.Email).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	if user.Confirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
		return
	}

	user.Email = input.Email
	err = user.SetPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting password"})
		return
	}

	confirmationCode := GenerateConfirmationCode()
	user.ConfirmedCode = confirmationCode
	user.ConfirmationTries = MaxConfirmationTries

	result := db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	emailBody := fmt.Sprintf("Welcome to our app! Your confirmation code is: %s", confirmationCode)
	err = SendEmail(user.Email, "Confirmation Code", emailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending confirmation code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation code sent to email"})
}

func ConfirmRegistration(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing confirmation code"})
		return
	}

	var user User
	err := db.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if input.Code != user.ConfirmedCode {
		user.ConfirmationTries--
		if user.ConfirmationTries == 0 {
			confirmationCode := GenerateConfirmationCode()
			user.ConfirmedCode = confirmationCode
			user.ConfirmationTries = MaxConfirmationTries

			emailBody := fmt.Sprintf("Your new confirmation code is: %s", confirmationCode)
			err = SendEmail(user.Email, "New Confirmation Code", emailBody)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending new confirmation code"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect confirmation code. A new code has been sent to your email."})
		} else {
			db.Save(&user)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Incorrect confirmation code. %v tries left.", user.ConfirmationTries)})
		}
		return
	}

	user.Confirmed = true
	db.Save(&user)

	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User confirmed successfully", "token": token})
}

func LoginUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if !user.Confirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not confirmed"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "token": token})
}

func ResendConfirmationCode(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	confirmationCode := GenerateConfirmationCode()
	user.ConfirmedCode = confirmationCode
	user.ConfirmationTries = MaxConfirmationTries

	result := db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	emailBody := fmt.Sprintf("Welcome to our app! Your confirmation code is: %s", confirmationCode)
	err := SendEmail(user.Email, "Confirmation Code", emailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending confirmation code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmation code sent to email"})
}

func SendEmail(to string, subject string, body string) error {
	from := smtpUser
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort),
		smtp.PlainAuth("", smtpUser, smtpPass, smtpServer),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return err
	}

	fmt.Println("Mail sent successfully!")
	return nil
}
