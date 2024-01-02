package main

import (
	"errors"
	"strings"
)

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Id       int    `json:"id"`
	Body     string `json:"body"`
}

func (db *DB) CreateChirp(body string, authId int) (Chirp, error) {
	if err := db.ensureDB(); err != nil {
		return Chirp{}, err
	}
	dbBody, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	if len(body) > 141 {
		return Chirp{}, errors.New("chirp is too long")
	}

	index := len(dbBody.Chirps) + 1
	body = removeProfane(body)
	return Chirp{
		AuthorId: authId,
		Id:       index,
		Body:     body,
	}, nil
}

func removeProfane(msg string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
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
	return respText
}
