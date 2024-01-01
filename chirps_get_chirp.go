package main

func (db *DB) GetChirps() ([]Chirp, error) {
	if err := db.ensureDB(); err != nil {
		return []Chirp{}, err
	}

	dbBody, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	allChirps := []Chirp{}
	for _, chrp := range dbBody.Chirps {
		allChirps = append(allChirps, chrp)
	}

	return allChirps, nil
}
