package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (db *DB) usersPostHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
	}

	usr, err := db.createUser(params.Body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
	}
	dbstruct := DBStructure{}
	if dbstruct.Users == nil {
		dbstruct.Users = map[int]User{}
	}
	dbstruct.Users[usr.Id] = usr
	db.writeDB(dbstruct)
	w.WriteHeader(201)
	respondWithJson(w, usr)
}
