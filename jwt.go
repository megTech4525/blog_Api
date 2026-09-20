package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func loadJwtSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT env variable is not set")
		return
	}
	jwtSecret = []byte(secret)
}

func generateToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"User_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
func validateToken(tokenString string) (int, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	rawID, exists := claims["User_id"]
	if !exists {
		return 0, fmt.Errorf("user_id claim missing")
	}
	idFloat, ok := rawID.(float64)
	if !ok {
		return 0, fmt.Errorf("User_id claim has wrong type")
	}
	return int(idFloat), nil
}
