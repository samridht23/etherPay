package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

func CorsMiddleware() func(http.Handler) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Replace with your frontend URL, or use "*" to allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           4000, // Maximum value not ignored by any of major browsers
	}
	return cors.New(corsOptions).Handler
}
