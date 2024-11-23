package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Get root path of project
var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b), "../../")
)

// Load env file to system.
func LoadEnv() {

	// load .env file from project root
	if err := godotenv.Load(Root + "/.env"); err != nil {
		log.Printf("Config: %v", err)
	}
}

// Get env value from key, default to second param
func GetEnvWithKey(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}
