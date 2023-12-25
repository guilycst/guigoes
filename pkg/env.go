package pkg

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	POSTS_PATH     string
	DIST_PATH      string
	BLEVE_IDX_PATH string
	SMTP_ENDPOINT  string
	SMTP_PORT      int
	SMTP_USR_NAME  string
	SMTP_USR_PW    string
)

func LoadEnvFile(filenames ...string) {
	godotenv.Load(filenames...)
	LoadEnvFromOS()
}

func LoadEnvFromOS() {
	POSTS_PATH = os.Getenv("POSTS_PATH")
	DIST_PATH = os.Getenv("DIST_PATH")
	BLEVE_IDX_PATH = os.Getenv("BLEVE_IDX_PATH")
	SMTP_ENDPOINT = os.Getenv("SMTP_ENDPOINT")

	var err error
	SMTP_PORT, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println(err)
	}

	SMTP_USR_NAME = os.Getenv("SMTP_USR_NAME")
	SMTP_USR_PW = os.Getenv("SMTP_USR_PW")
}
