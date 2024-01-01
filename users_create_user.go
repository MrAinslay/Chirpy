package main

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) createUser(email string) (User, error) {
	if err := db.ensureDB(); err != nil {
		return User{}, err
	}

	datbase, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	index := len(datbase.Users) + 1
	return User{
		Id:    index,
		Email: email,
	}, nil
}
