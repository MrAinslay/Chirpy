package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type responseError struct {
		Error string `json:"error"`
	}
	respErr := responseError{
		Error: msg,
	}
	w.WriteHeader(code)
	dat, _ := json.Marshal(respErr)
	w.Write(dat)
}
