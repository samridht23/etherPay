package main

import (
	"os"

	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/configs"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/infra"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/portal"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/router"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Application struct {
	Router *chi.Mux
}

func InitializeApplication() *Application {
	return &Application{
		Router: chi.NewRouter(),
	}
}

func (app *Application) New() error {

	// initialize enviroment variable
	configs.InitializeEnvVariables()

	// initialize a custom configured logger
	infra.InitializeZapLogger()

	dbPool, err := infra.InitializePool()
	if err != nil {
		return err
	}

	pool := dbPool.New()
	defer pool.Close()

	// initialize aws s3 client
	s3, err := infra.InitializeS3Client()
	if err != nil {
		return err
	}
	s3Client := s3.NewClient()

	router.InitializeRoutes(app.Router, s3Client, pool)
	portal.InitializeServer(app.Router)
	return nil
}

func main() {
	app := InitializeApplication()
	err := app.New()
	if err != nil {
		zap.S().Errorf("Error initalizing application: %v", err)
		os.Exit(1)
	}
}
