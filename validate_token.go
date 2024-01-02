package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func validateToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	bearerTkn := r.Header.Get("Authorization")

	slpitTkn := strings.Split(bearerTkn, " ")
	tkn := slpitTkn[1]

	godotenv.Load("key.env")
	jwtKey := os.Getenv("JWT_SECRET")

	jwtTkn, err := jwt.ParseWithClaims(tkn, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	}

	return jwtTkn, nil
}
