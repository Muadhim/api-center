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
	SmptHost        string
	SmptPort        string
	SmptUser        string
	SmptPass        string
}

func LoadConfig() *Config {
	config := &Config{}
	envVars := map[string]*string{
		"POSTGRES_HOST":      &config.DBHost,
		"DB_PORT":            &config.DBPort,
		"POSTGRES_USER":      &config.DBUser,
		"POSTGRES_PASSWORD":  &config.DBPass,
		"POSTGRES_DATABASE":  &config.DBName,
		"SSL_MODE":           &config.SSLMode,
		"SERVER_PORT":        &config.ServerPort,
		"SET_MAX_IDLE_CONNS": &config.SetMaxIdleConns,
		"SET_MAX_OPEN_CONNS": &config.SetMaxOpenConns,
		"API_VERSION":        &config.ApiVersion,
		"SMTP_HOST":          &config.SmptHost,
		"SMTP_PORT":          &config.SmptPort,
		"SMTP_AUTH_EMAIL":    &config.SmptUser,
		"SMTP_AUTH_PASSWORD": &config.SmptPass,
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
