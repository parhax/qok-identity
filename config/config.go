package config

import (
	"os"
)

type Config struct {
	Mongo_url string
	Http_port string
}

func Load() *Config {
	return &Config{
		Mongo_url: getEnv("MONGO_URL", "mongodb://localhost:27017"),
		Http_port: getEnv("HTTP_PORT", "8585"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
