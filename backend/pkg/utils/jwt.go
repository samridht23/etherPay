package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaim represents custom JWT claims including registered claims
type CustomJWTClaim struct {
	EncryptedValue string `json:"encrypted_value"`
	jwt.RegisteredClaims
}

// generate signed JWT token  with given value and key
func SignJWTToken(value string, key []byte) (string, error) {

	// create custom claim
	claim := CustomJWTClaim{
		EncryptedValue: value,
	}

	// create a new token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// sign the token using provided key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJWTToken(tk string, key []byte) (*CustomJWTClaim, error) {

	token, err := jwt.ParseWithClaims(tk, &CustomJWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid jwt signing method detected while parsing: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("JWT token invalid")
	}

	claims, ok := token.Claims.(*CustomJWTClaim)
	if !ok {
		return nil, fmt.Errorf("Error parsing jwt claim")
	}

	return claims, nil
}
