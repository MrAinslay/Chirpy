package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (db *DB) getHandler(w http.ResponseWriter, r *http.Request) {
	allChirps, err := db.GetChirps()
	strId := r.URL.Query().Get("author_id")
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if strId == "" {
		for _, chrp := range allChirps {
			respondWithJson(w, chrp)
		}
		return
	} else {
		id, err := strconv.Atoi(strId)
		if err != nil {
			respondWithError(w, 500, fmt.Sprint(err))
			return
		} else {
			for _, chrp := range allChirps {
				if chrp.AuthorId == id {
					respondWithJson(w, chrp)
				}
			}
			return
		}
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
