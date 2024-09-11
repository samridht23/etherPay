package middleware

import (
	"context"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/utils"
	"net/http"
	"os"
)

func AuthorizeMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("_acc_tk")
			if err != nil {
				if err == http.ErrNoCookie {
					utils.Err(w, http.StatusBadRequest, []string{"Unauthorized request"})
					return
				}
				utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
				return
			}

			jwtClaim, err := utils.ParseJWTToken(cookie.Value, []byte(os.Getenv("JWT_KEY")))
			if err != nil {
				utils.Err(w, http.StatusUnauthorized, []string{"Unauthorized request"})
				return
			}

			address := jwtClaim.EncryptedValue

			userCtxValue := &pkg.UserContextData{
				Address: address,
			}

			ctx := context.WithValue(r.Context(), "auth_context", userCtxValue)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
