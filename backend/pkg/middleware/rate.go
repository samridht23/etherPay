package middleware

import (
	"net/http"
	"time"

	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/utils"
	"golang.org/x/time/rate"
)

func RateLimiter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := rate.NewLimiter(rate.Every(time.Second), 5) // 5 burst request
			if limiter != nil {
				if !limiter.Allow() {
					utils.Err(w, http.StatusTooManyRequests, []string{"Request rate limit exceeded. Try after some time"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
