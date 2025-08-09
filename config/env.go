package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DBSSLMode     string
	DBSSLRootCert string
}

func LoadConfig() *Config {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found")
		}
	}

	return &Config{
		Port:          getEnv("PORT"),
		DBUser:        getEnv("DB_USER"),
		DBPassword:    getEnv("DB_PASSWORD"),
		DBHost:        getEnv("DB_HOST"),
		DBPort:        getEnv("DB_PORT"),
		DBName:        getEnv("DB_NAME"),
		DBSSLMode:     getEnv("DB_SSL_MODE"),
		DBSSLRootCert: getEnv("DB_SSL_ROOT_CERT"),
	}
}

func getEnv(name string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable not found: %v", name))
}
