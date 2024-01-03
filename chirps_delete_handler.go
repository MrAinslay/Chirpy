package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (db *DB) chirpsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	jwtTkn, err := validateToken(w, r)
	if err != nil {
		respondWithError(w, 401, fmt.Sprint(err))
		return
	}

	strId := chi.URLParam(r, "id")
	if strId == "" {
		respondWithError(w, 404, "page not found")
		return
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	datbase, _ := db.loadDB()

	subjStrbjId, _ := jwtTkn.Claims.GetSubject()
	subjId, _ := strconv.Atoi(subjStrbjId)
	if datbase.Chirps[id].AuthorId != subjId {
		respondWithError(w, 403, "")
		return
	}

	delete(datbase.Chirps, id)
	log.Println(datbase.Chirps)
	db.deleteData("chrip", id)
}
