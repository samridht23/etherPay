package configs

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// We will use the docker-compose enviroment parameter to pass all the
// enviroment and export it to the container shell running the server image
// eg. APP_ENV or SERVER_PORT

// In production - the env variables will get loaded from the shell
// enviroment in which the application is running (docker container)

// In development - the env variables will be loaded from .env file
// also in development will be pass the first env variable through Makefile
// which is used to detect what enviroment the server is running (prod or dev)

// InitializeEnvVariables loads environment variables based on the application environment.
func InitializeEnvVariables() {
	env := os.Getenv("APP_ENV")

	switch env {
	case "development":
		loadDevelopmentEnv()
	case "production":
		slog.Info("Production environment variables loaded.")
	case "":
		slog.Error("No application environment provided. Exiting.")
		os.Exit(1)
	default:
		slog.Error("Invalid application environment provided: %s. Exiting.", env)
		os.Exit(1)
	}
}

// loadDevelopmentEnv loads environment variables from a .env file for development purposes.
func loadDevelopmentEnv() {
	slog.Info("Loading development environment variables...")
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading environment variables: %v", err)
		os.Exit(1)
	}
	slog.Info("Development environment variables loaded successfully.")
}
