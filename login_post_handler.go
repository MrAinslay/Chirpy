package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) loginHnalder(w http.ResponseWriter, r *http.Request) {
	type loginInfo struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type loginResp struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	loginInf := loginInfo{}

	if err := decoder.Decode(&loginInf); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		return
	}

	if err := godotenv.Load("key.env"); err != nil {
		respondWithError(w, 400, fmt.Sprint(err))
		log.Println("wot", err)
		return
	}
	key := os.Getenv("JWT_SECRET")

	cfg := apiConfig{
		jwtKey: key,
	}

	datbase, _ := db.loadDB()

	for _, usr := range datbase.Users {
		if usr.Email == loginInf.Email {
			err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(loginInf.Password))
			if err != nil {
				respondWithError(w, 401, fmt.Sprint(err))
				return
			}

			tkn := cfg.createToken("chirpy-access", usr.Id)
			refreshTkn := cfg.createToken("chirpy-refresh", usr.Id)
			loginRsp := loginResp{
				Id:           usr.Id,
				Email:        usr.Email,
				AccessToken:  tkn,
				RefreshToken: refreshTkn,
			}
			respondWithJson(w, loginRsp)
			return
		}
	}
	respondWithError(w, 400, "Could not find user with this email")
}
