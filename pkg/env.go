package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

var POSTS_PATH string

func LoadEnvFile() {
	godotenv.Load()
	LoadEnvFromOS()
}

func LoadEnvFromOS() {
	POSTS_PATH = os.Getenv("POSTS_PATH")
}
