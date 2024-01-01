package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, data any) {
	dat, _ := json.Marshal(data)
	w.Write(dat)
}
