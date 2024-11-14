package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectionStringDB string
	Port               int
	SECRET_KEY         []byte
)

// Initialize reads all the environment variables.
func Initialize() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Without the env vars we cannot initialize the API server.\n", err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))

	if err != nil {
		Port = 9000
	}

	ConnectionStringDB = fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
}
