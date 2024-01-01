package main

import (
	"net/http"
)

func (db *DB) getHandler(w http.ResponseWriter, r *http.Request) {
	allChirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	w.WriteHeader(200)
	for _, chrp := range allChirps {
		respondWithJson(w, chrp)
	}
}
