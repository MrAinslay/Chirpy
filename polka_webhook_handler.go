package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (db *DB) polkaWebhookHanlder(w http.ResponseWriter, r *http.Request) {
	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	authHeaderKey := r.Header.Get("Authorization")
	rqstKey := strings.Split(authHeaderKey, " ")

	if len(rqstKey) <= 1 {
		w.WriteHeader(401)
		return
	}

	if rqstKey[1] != db.cfg.apiKey {
		w.WriteHeader(401)
		return
	}

	docoder := json.NewDecoder(r.Body)
	rqst := polkaRequest{}
	if err := docoder.Decode(&rqst); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	if rqst.Event != "user.upgraded" {
		w.WriteHeader(200)
	}

	datbase, _ := db.loadDB()

	for _, usr := range datbase.Users {
		if rqst.Data.UserId == usr.Id {
			usr.ChirpyRed = true
			datbase.Users[usr.Id] = usr
			db.writeDB(datbase)
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(404)
}
