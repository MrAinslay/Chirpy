package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (db *DB) getIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(400)
		return
	}
	allChirps, err := db.GetChirps()
	if err != nil {
		w.WriteHeader(400)
		return
	}

	for _, chrp := range allChirps {
		if fmt.Sprint(chrp.Id) == id {
			dat, err := json.Marshal(chrp)
			if err != nil {
				w.WriteHeader(400)
			}
			w.Write(dat)
			return
		}
	}
	w.WriteHeader(404)
}
