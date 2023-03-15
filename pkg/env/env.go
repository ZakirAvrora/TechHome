package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const env = "ENV"
const local = "local"

func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" && key != "PORT" {
		panic("env key " + val + "cannot found")
	}
	return val
}

func CheckDotEnv(path string) {
	err := godotenv.Load(path)
	if err != nil && os.Getenv(env) == local {
		log.Println("Error in loading .env file:", err)
	}
}
