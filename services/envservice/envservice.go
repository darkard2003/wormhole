package envservice

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	os.Setenv("ENV_LOADED", "true")
	log.Println("Environment variables loaded successfully")
	return nil
}

func GetEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", key)
	}
	return value, nil
}
