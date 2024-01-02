package main

import "net/http"

func (db *DB) putHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
