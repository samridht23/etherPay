package router

import (
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/controller"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/middleware"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func InitializeRoutes(router *chi.Mux, s3Client *s3.Client, pool *pgxpool.Pool) {
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.RateLimiter())
	router.Use(chiMiddleware.Heartbeat("/"))

	router.Group(func(router chi.Router) {
		router.Use(middleware.AuthorizeMiddleware())
		router.Get("/auth-status", controller.AuthStatus(pool))
		router.Post("/transaction", controller.Transaction(pool))
		router.Post("/profile", controller.UpdateUserData(pool))
	})

	router.Get("/profile/{address}", controller.GetUser(pool))
	router.Post("/connect", controller.Connect(pool))
}
