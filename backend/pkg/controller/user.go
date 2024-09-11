package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/utils"
	//"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Address           string      `json:"address"`
	UserName          pgtype.Text `json:"username"`
	ProfilePictureUrl pgtype.Text `json:"profile_image_url"`
	BannerImageUrl    pgtype.Text `json:"banner_image_url"`
	About             pgtype.Text `json:"about"`
	CreatedAt         time.Time   `json:"created_at"`
}

func GetUser(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAddress := chi.URLParam(r, "address")
		var user User
		query := `SELECT * FROM "User" WHERE "address" = $1`
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := pool.QueryRow(ctx, query, userAddress).Scan(
			&user.Address,
			&user.UserName,
			&user.ProfilePictureUrl,
			&user.BannerImageUrl,
			&user.About,
			&user.CreatedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				utils.Err(w, http.StatusNotFound, []string{"User not found"})
				return
			}
			log.Println(err)
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		log.Printf("User data asked by the client %+v\n", user)
		utils.Send(w, http.StatusOK, user)
		return
	}
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	About    string `json:"about"`
}

func UpdateUserData(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve authorization context
		ctxValue, ok := r.Context().Value("auth_context").(*pkg.UserContextData)
		if !ok {
			log.Println("Authorization context value not found or invalid")
			utils.Err(w, http.StatusUnauthorized, []string{"Unauthorized request"})
			return
		}

		// Decode JSON request body
		var updateRequest UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			utils.Err(w, http.StatusBadRequest, []string{"Failed to parse request body"})
			return
		}

		// Check for empty fields and prepare update query
		var query string
		var args []interface{}

		if updateRequest.Username == "" && updateRequest.About == "" {
			utils.Err(w, http.StatusBadRequest, []string{"Missing required fields"})
			return
		} else if updateRequest.Username == "" {
			query = `UPDATE "User" SET "about" = $1 WHERE "address" = $2`
			args = []interface{}{updateRequest.About, ctxValue.Address}
		} else if updateRequest.About == "" {
			query = `UPDATE "User" SET "username" = $1 WHERE "address" = $2`
			args = []interface{}{updateRequest.Username, ctxValue.Address}
		} else {
			query = `UPDATE "User" SET "username" = $1, "about" = $2 WHERE "address" = $3`
			args = []interface{}{updateRequest.Username, updateRequest.About, ctxValue.Address}
		}

		// Execute update query
		_, err := pool.Exec(r.Context(), query, args...)
		if err != nil {
			log.Println("Failed to update user data:", err)
			utils.Err(w, http.StatusInternalServerError, []string{"Failed to update user data"})
			return
		}

		// Send success response
		utils.Send(w, http.StatusOK, []string{"User updated successfully"})
	}
}
