package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func loadEnv() Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}

	return Config{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
}
