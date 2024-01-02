package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) putHandler(w http.ResponseWriter, r *http.Request) {
	type respBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	rsp := respBody{}

	if err := decoder.Decode(&rsp); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	bearerTkn := r.Header.Get("Authorization")

	slpitTkn := strings.Split(bearerTkn, " ")
	tkn := slpitTkn[1]

	godotenv.Load("key.env")
	jwtKey := os.Getenv("JWT_SECRET")

	jwtTkn, err := jwt.ParseWithClaims(tkn, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		respondWithError(w, 401, fmt.Sprint(err))
		return
	}

	if isr, _ := jwtTkn.Claims.GetIssuer(); isr == "chirpy-refresh" {
		respondWithError(w, 401, "does not accept refresh tokens")
		return
	}

	id, err := jwtTkn.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	encrPass, err := bcrypt.GenerateFromPassword([]byte(rsp.Password), 10)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	err = db.updateUser(id, rsp.Email, string(encrPass))
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}
	w.WriteHeader(200)
	respondWithJson(w, rsp)
}
