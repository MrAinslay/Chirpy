package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) loginHnalder(w http.ResponseWriter, r *http.Request) {
	type loginInfo struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type loginResp struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	loginInf := loginInfo{}

	if err := decoder.Decode(&loginInf); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	datbase, _ := db.loadDB()

	for _, usr := range datbase.Users {
		if usr.Email == loginInf.Email {
			err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(loginInf.Password))
			if err != nil {
				respondWithError(w, 401, fmt.Sprint(err))
				return
			}
			loginRsp := loginResp{
				Id:    usr.Id,
				Email: usr.Email,
			}
			respondWithJson(w, loginRsp)
			return
		}
	}
	respondWithError(w, 400, "Could not find user with this email")
}
