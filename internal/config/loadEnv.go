package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// LoadEnvVariables loads the .env file from project root or relative paths.
func LoadEnvVariables() {
	// List of possible .env locations relative to the working directory
	possiblePaths := []string{
		".env",          // current directory
		"../.env",       // parent
		"../../.env",    // grandparent
		"../../../.env", // etc
	}

	var loaded bool
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err != nil {
				log.Fatalf("Error loading .env file from %s: %v", p, err)
			}
			loaded = true
			log.Printf("Loaded .env from %s", p)
			break
		}
	}

	if !loaded {
		log.Println("No .env file found in expected locations. Continuing without it...")
	}

	// Set Gin mode from environment variable, fallback to debug
	mode := os.Getenv("GIN_MODE")
	switch mode {
	case gin.ReleaseMode, gin.DebugMode, gin.TestMode:
		gin.SetMode(mode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
