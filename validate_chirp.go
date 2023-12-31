package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type error struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if len(params.Body) > 140 {
		w.WriteHeader(400)
		respBody := error{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		respBody := error{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
		return
	}
	w.WriteHeader(http.StatusOK)

	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(params.Body, " ")
	for i, word := range words {
		for _, profane := range profanity {
			if strings.ToLower(word) == profane {
				words[i] = "****"
			}
		}
	}

	respText := ""
	for _, word := range words {
		respText += word + " "
	}
	respText = strings.Trim(respText, " ")
	respBody := parameters{
		Body: respText,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}
