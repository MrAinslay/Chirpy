package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps         map[int]Chirp        `json:"chirps"`
	Users          map[int]User         `json:"users"`
	RevockedTokens map[int]RevokedToken `json:"revoked_tokens"`
}

func NewDB(path string) (*DB, error) {
	databaseData := []byte(`{
		"chirps": {
		}
	  }`)
	err := os.WriteFile("database.json", databaseData, 0666)
	if err != nil {
		log.Fatal("Could not create new database")
		return &DB{}, err
	}

	mux := &sync.RWMutex{}
	database := DB{
		path: path,
		mux:  mux,
	}
	return &database, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile("database.json")
	if err != nil {
		newDatabase, err2 := NewDB("database.json")
		if err2 != nil {
			return err2
		}
		db.path = newDatabase.path
		db.mux = newDatabase.mux
		return nil
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	dat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	databaseBody := DBStructure{}
	if err := json.Unmarshal(dat, &databaseBody); err != nil {
		return DBStructure{}, err
	}
	return databaseBody, nil
}

func (db *DB) writeDB(dbstructure DBStructure) error {
	if err := db.ensureDB(); err != nil {
		return err
	}

	datbase, err := db.loadDB()
	if err != nil {
		return err
	}

	if datbase.Users == nil {
		datbase.Users = map[int]User{}
	}
	if datbase.RevockedTokens == nil {
		datbase.RevockedTokens = map[int]RevokedToken{}
	}

	for index, chrp := range dbstructure.Chirps {
		datbase.Chirps[index] = chrp
	}
	for index, usr := range dbstructure.Users {
		datbase.Users[index] = usr
	}
	for index, tkn := range dbstructure.RevockedTokens {
		datbase.RevockedTokens[index] = tkn
	}

	dat, err2 := json.Marshal(datbase)
	if err2 != nil {
		return err
	}
	defer db.mux.Unlock()
	db.mux.Lock()
	os.WriteFile("database.json", []byte(dat), 0666)
	return nil
}

func (db *DB) updateUser(id string, email string, pass string) error {
	datbase, err := db.loadDB()
	if err != nil {
		return err
	}

	usrId, _ := strconv.Atoi(id)
	if usr, ok := datbase.Users[usrId]; ok {
		usr.Email = email
		usr.Password = pass

		datbase.Users[usrId] = usr
	}

	err = db.writeDB(datbase)
	if err != nil {
		return err
	}
	return nil
}
