package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RevokedToken struct {
	Id        string `json:"id"`
	RevokedAt string `json:"revoked_at"`
	Token     string `json:"refresh_token"`
}

func (db *DB) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	bearerTkn := r.Header.Get("Authorization")
	if bearerTkn == "" {
		respondWithError(w, 400, "authorization header is required")
		return
	}

	strTkn := strings.Split(bearerTkn, " ")[1]

	datbase, _ := db.loadDB()
	if datbase.RevockedTokens == nil {
		datbase.RevockedTokens = map[int]RevokedToken{}
	}

	index := len(datbase.RevockedTokens) + 1
	datbase.RevockedTokens[index] = RevokedToken{
		Id:        strTkn,
		RevokedAt: fmt.Sprint(time.Now().UTC()),
		Token:     strTkn,
	}

	err := db.writeDB(datbase)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}
	w.WriteHeader(200)
}
