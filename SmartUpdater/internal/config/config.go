package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	GithubToken string
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  int
}

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

	if config.GithubToken == "" {
		logrus.Fatal("GITHUB_TOKEN environment variable is required")
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		logrus.Warnf("Invalid integer value for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
} 