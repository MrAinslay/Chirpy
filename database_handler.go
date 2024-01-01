package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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
	if err == os.ErrNotExist {
		newDatabase, erro := NewDB("database.json")
		db.path = newDatabase.path
		db.mux = newDatabase.mux
		return erro
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
	for index, chrp := range dbstructure.Chirps {
		datbase.Chirps[index-1] = chrp
	}
	return nil
}
