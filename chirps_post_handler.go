package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (db *DB) postHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	jwtTkn, err := validateToken(w, r)
	if err != nil {
		respondWithError(w, 401, fmt.Sprint(err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	strId, _ := jwtTkn.Claims.GetSubject()
	id, _ := strconv.Atoi(strId)

	chrp, err := db.CreateChirp(params.Body, id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
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
