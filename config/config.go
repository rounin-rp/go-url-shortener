package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DatabaseUrl  = ""
	DatabaseName = ""
)

func InitializeEnv() (bool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	DatabaseUrl = os.Getenv("DBURI")
	DatabaseName = os.Getenv("DBNAME")
	return true, nil
}
