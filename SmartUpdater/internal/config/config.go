package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	GithubToken string // GitHub API token
	DBHost      string // Database host
	DBPort      int    // Database port
	DBUser      string // Database user
	DBPassword  string // Database password
	DBName      string // Database name
	ServerPort  int    // HTTP server port
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		GithubToken: getEnvOrDefault("GITHUB_TOKEN", ""),
		DBHost:      getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:      getEnvAsIntOrDefault("DB_PORT", 5432),
		DBUser:      getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:  getEnvOrDefault("DB_PASSWORD", ""),
		DBName:      getEnvOrDefault("DB_NAME", "smartupdater"),
		ServerPort:  getEnvAsIntOrDefault("PORT", 8080),
	}

	// Validate required configuration
	if config.GithubToken == "" {
		logrus.Fatal("GITHUB_TOKEN environment variable is required")
	}

	return config
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault gets an environment variable as an integer or returns a default value
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		logrus.Warnf("Invalid integer value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
} 