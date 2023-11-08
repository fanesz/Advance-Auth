package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading .env file")
	}
	return os.Getenv(key)
}
