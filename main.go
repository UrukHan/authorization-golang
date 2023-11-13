package main

import (
	"auth/app"
	"fmt"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	// Setup configuration
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	app.MaxConfirmationTries = viper.GetInt("maxConfirmationTries")

	// Load from Kubernetes secrets
	app.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	app.AdminEmail = os.Getenv("ADMIN_EMAIL")
	app.JwtSecret = os.Getenv("JWTSECRET")

	// Connect to PostgreSQL
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_DBNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SSLMODE")))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&app.User{})

	// Adding routes to the same router
	r := app.SetupRoutes(db)

	// Run un 0.0.0.0:8020
	r.Run(":8020")
}