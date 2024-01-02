package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) usersPostHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresAt string `json:"expires_in_seconds"`
	}
	type response struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	encrPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	usr, err := db.createUser(params.Email, string(encrPass), params.ExpiresAt)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}
	dbstruct := DBStructure{}
	if dbstruct.Users == nil {
		dbstruct.Users = map[int]User{}
	}
	dbstruct.Users[usr.Id] = usr
	db.writeDB(dbstruct)
	w.WriteHeader(201)

	rsp := response{
		Id:    usr.Id,
		Email: usr.Email,
	}
	respondWithJson(w, rsp)
}
