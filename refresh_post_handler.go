package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func (db *DB) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type respBody struct {
		Token string `json:"token"`
	}

	if err := godotenv.Load("key.env"); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}
	jwtKey := os.Getenv("JWT_SECRET")
	cfg := apiConfig{
		jwtKey: jwtKey,
	}

	bearerTkn := r.Header.Get("Authorization")
	if bearerTkn == "" {
		respondWithError(w, 400, "authorization header is required")
		return
	}

	strTkn := strings.Split(bearerTkn, " ")[1]

	jwtTkn, err := jwt.ParseWithClaims(strTkn, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		respondWithError(w, 401, fmt.Sprint(err))
		return
	}

	if isr, _ := jwtTkn.Claims.GetIssuer(); isr != "chirpy-refresh" {
		respondWithError(w, 401, "not a refresh token")
		return
	}

	datbase, _ := db.loadDB()
	for _, tkn := range datbase.RevockedTokens {
		if tkn.Token == strTkn {
			respondWithError(w, 401, "refresh token is revoked")
			return
		}
	}

	strId, _ := jwtTkn.Claims.GetSubject()
	id, _ := strconv.Atoi(strId)

	rsp := respBody{
		Token: cfg.createToken("chirpy-access", id),
	}

	respondWithJson(w, rsp)
}
