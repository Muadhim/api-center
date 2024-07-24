package api

import (
	"api-center/api/controllers"
	"api-center/configs"
	"log/slog"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	var err error = godotenv.Load()

	if err != nil {
		slog.Warn("Error loading .env file")
	}
	// Load configuration
	config := configs.LoadConfig()

	server.Initialize(config)

	// seeds.Load(server.DB)

	server.Run(config.ServerPort)
}
