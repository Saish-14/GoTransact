package config

import (
	"log"
	"os"
	// "fmt"
	"github.com/joho/godotenv"
)

var (
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbTimezone string
	FromEmail string
	EmailUser string
	EmailPass string
	EmailHost string
	EmailPort string
)

func Loadenv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file ERROR:", err)
	}
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbName = os.Getenv("DB_NAME")
	DbTimezone = os.Getenv("DB_TIMEZONE")
	FromEmail = os.Getenv("FROM_EMAIL")
	EmailUser = os.Getenv("EMAIL_USER")
	EmailPass = os.Getenv("EMAIL_PASS")
	EmailHost = os.Getenv("EMAIL_HOST")
	EmailPort = os.Getenv("EMAIL_PORT")
}

//dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbTimezone)

