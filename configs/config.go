package configs

import (
	"log"
	"os"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	SSLMode         string
	ServerPort      string
	SetMaxIdleConns string
	SetMaxOpenConns string
	ApiVersion      string
}

func LoadConfig() *Config {
	config := &Config{}
	envVars := map[string]*string{
		"DB_HOST":            &config.DBHost,
		"DB_PORT":            &config.DBPort,
		"DB_USER":            &config.DBUser,
		"DB_PASS":            &config.DBPass,
		"DB_NAME":            &config.DBName,
		"SSL_MODE":           &config.SSLMode,
		"SERVER_PORT":        &config.ServerPort,
		"SET_MAX_IDLE_CONNS": &config.SetMaxIdleConns,
		"SET_MAX_OPEN_CONNS": &config.SetMaxOpenConns,
		"API_VERSION":        &config.ApiVersion,
	}

	for key, ptr := range envVars {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("Missing environment variable: %s", key)
		}
		*ptr = value
	}

	return config
}
