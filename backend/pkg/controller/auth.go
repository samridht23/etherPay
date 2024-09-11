package controller

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg"
	"github.com/W0NB0N/buymeacrypto-demo-/backend/pkg/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

func verifySignature(address, signature, message string) (bool, error) {
	msgHash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	// Decode the signature, strip the "0x" prefix
	sigBytes, err := hex.DecodeString(signature[2:])
	if err != nil {
		return false, err
	}

	// Ensure the signature is of the correct length (65 bytes)
	if len(sigBytes) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// Adjust the 'v' value
	if sigBytes[64] != 27 && sigBytes[64] != 28 {
		return false, fmt.Errorf("invalid signature v byte")
	}
	sigBytes[64] -= 27

	// Recover the public key
	pubKey, err := crypto.SigToPub(msgHash.Bytes(), sigBytes)
	if err != nil {
		return false, err
	}

	// Recover the Ethereum address from the public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	return strings.ToLower(recoveredAddress) == strings.ToLower(address), nil
}

func CheckUserExists(address string, pool *pgxpool.Pool) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "User" WHERE address = $1);`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := pool.QueryRow(ctx, query, address).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("Query failed to check user exists: %w", err)
	}
	return exists, nil
}

func AddNewUser(address string, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("Could not begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background()) // Rollback if function exits without commit
	newUserQuery := `
		INSERT INTO "User" (address,created_at)
		VALUES ($1, $2)
	`
	userData := []interface{}{
		address,
		time.Now().UTC(),
	}
	_, err = tx.Exec(context.Background(), newUserQuery, userData...)
	if err != nil {
		return fmt.Errorf("error inserting new user: %v", err)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("Error committing transaction: %v", err)
	}
	return nil
}

func Connect(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.Err(w, http.StatusBadRequest, []string{"Invalid request body"})
			return
		}
		if len(req.Address) == 0 {
			utils.Err(w, http.StatusBadRequest, []string{"Address is required"})
			return
		}
		log.Println(req.Address)
		log.Println(req.Signature)
		log.Println(req.Message)

		isValid, err := verifySignature(req.Address, req.Signature, req.Message)
		if err != nil {
			log.Println("triggered here")
			log.Println(err)
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		log.Println(isValid)
		if !isValid {
			utils.Err(w, http.StatusUnauthorized, []string{"Invalid signature"})
			return
		}
		userExists, err := CheckUserExists(req.Address, pool)
		if err != nil {
			log.Println(err)
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		log.Println("user existence", userExists)
		if userExists {
			signedToken, err := utils.SignJWTToken(req.Address, []byte(os.Getenv("JWT_KEY")))
			if err != nil {
				log.Println(err)
				utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
				return
			}
			cookie := http.Cookie{
				Name:     "_acc_tk",
				Value:    signedToken,                         // You can use any generated value like JWT
				Expires:  time.Now().Add(24 * time.Hour * 30), // Set expiration time for the cookie
				Path:     "/",
				HttpOnly: true,                 // Ensures the cookie is not accessible via JavaScript
				Secure:   false,                // Set true if using HTTPS
				SameSite: http.SameSiteLaxMode, // Lax mode is a good default for most use cases
			}
			log.Printf("%+v", &cookie)
			http.SetCookie(w, &cookie)
			utils.Send(w, http.StatusOK, []string{"User logged in successfully"})
			return
		}

		err = AddNewUser(req.Address, pool)
		if err != nil {
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		signedToken, err := utils.SignJWTToken(req.Address, []byte(os.Getenv("JWT_KEY")))
		if err != nil {
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		cookie := http.Cookie{
			Name:     "_acc_tk",
			Value:    signedToken,
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour), // Set cookie to expire in 30 days
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
		redirectURL := os.Getenv("CLIENT_URL")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
}

type AuthStatusData struct {
	Address         string      `json:"address"`
	Username        pgtype.Text `json:"username"`
	ProfileImageUrl pgtype.Text `json:"profile_image_url"`
	BannerImageUrl  pgtype.Text `json:"banner_image_url"`
}

func AuthStatus(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctxValue, ok := r.Context().Value("auth_context").(*pkg.UserContextData)
		if !ok {
			log.Println("Authorization context value not found or invalid")
			utils.Err(w, http.StatusUnauthorized, []string{"Unauthorized request"})
			return
		}

		var authStatsuData AuthStatusData

		query := `
      SELECT
          u.address,
          u.username,
          u.profile_image_url,
          u.banner_image_url
      FROM
          "User" u
      WHERE
          u.address = $1;
    `

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := pool.QueryRow(ctx, query, ctxValue.Address).Scan(
			&authStatsuData.Address,
			&authStatsuData.Username,
			&authStatsuData.ProfileImageUrl,
			&authStatsuData.BannerImageUrl,
		)

		if err != nil {
			log.Println(err)
			utils.Err(w, http.StatusInternalServerError, []string{"Internal server error"})
			return
		}
		utils.Send(w, http.StatusOK, authStatsuData)
	}
}
