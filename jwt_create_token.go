package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) createToken(issuer string, id int) string {
	key := cfg.jwtKey
	expiration := 0
	if issuer == "chirpy-refresh" {
		expiration = 1440
	} else {
		expiration = 1
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  &jwt.NumericDate{time.Now().UTC()},
		ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Duration(expiration) * time.Hour)},
		Subject:   fmt.Sprint(id),
	})
	signedTkn, _ := tkn.SignedString([]byte(key))
	return signedTkn
}
