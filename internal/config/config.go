package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	Env        string
	JWTSecret  string
	BaseURL    string
}

var (
	config *Config
	once   sync.Once
)

// GetConfig returns the singleton instance of Config
func GetConfig() *Config {
	once.Do(func() {
		var err error
		config, err = loadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
	})
	return config
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "shortlink"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		BaseURL:    getEnv("BASE_URL", "http://localhost"),
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
