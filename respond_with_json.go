package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, chrp Chirp) {
	dat, _ := json.Marshal(chrp)
	w.Write(dat)
}
