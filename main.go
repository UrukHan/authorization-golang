package main

import (
	"auth/app"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	// Read configuration file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// Variables from config
	app.MaxConfirmationTries = viper.GetInt("maxConfirmationTries")
	app.AdminDataPath = viper.GetString("adminDataPath")
	app.JwtSecret = viper.GetString("jwtSecret")

	// Connect postgress
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.dbname"),
		viper.GetString("db.password"),
		viper.GetString("db.sslmode")))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&app.User{})

	// Adding routes to the same router
	r := app.SetupRoutes(db)

	// Run auto subscriber function in goroutine
	go app.AutoRenewSubscriptions(db)

	// Run un 0.0.0.0:8020
	r.Run(":8020")
}
