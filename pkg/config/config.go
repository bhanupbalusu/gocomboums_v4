package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
	LogLevel   string
}

func LoadConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432")) // Default port is 5432
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "user_management"),
		DBPort:     dbPort,
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}, nil
}

// getEnv gets the environment variable if present, else returns the default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
