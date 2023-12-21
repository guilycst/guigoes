package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	POSTS_PATH      string
	SERVICE         string
	GITHUB_API_HOST string
	REPO_OWNER      string
	POSTS_REPO      string
	DIST_PATH       string
	BLEVE_IDX_PATH  string
)

func LoadEnvFile() {
	godotenv.Load()
	LoadEnvFromOS()
}

func LoadEnvFromOS() {
	POSTS_PATH = os.Getenv("POSTS_PATH")
	SERVICE = os.Getenv("SERVICE")
	GITHUB_API_HOST = os.Getenv("GITHUB_API_HOST")
	REPO_OWNER = os.Getenv("REPO_OWNER")
	POSTS_REPO = os.Getenv("POSTS_REPO")
	DIST_PATH = os.Getenv("DIST_PATH")
	BLEVE_IDX_PATH = os.Getenv("BLEVE_IDX_PATH")
}
