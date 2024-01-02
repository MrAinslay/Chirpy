package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createToken(expiration string, id int) string {
	key := cfg.jwtKey
	seconds, _ := strconv.Atoi(expiration)
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{time.Now().UTC()},
		ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Duration(seconds))},
		Subject:   fmt.Sprint(id),
	})
	signedTkn, _ := tkn.SignedString([]byte(key))
	return signedTkn
}
