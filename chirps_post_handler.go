package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (db *DB) postHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
	}

	chrp, err := db.CreateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
	}
	dbstruct := DBStructure{}
	if dbstruct.Chirps == nil {
		dbstruct.Chirps = map[int]Chirp{}
	}
	dbstruct.Chirps[chrp.Id] = chrp
	db.writeDB(dbstruct)
	w.WriteHeader(201)
	respondWithJson(w, chrp)
}
