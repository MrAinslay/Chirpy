package main

import "errors"

type User struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (db *DB) createUser(email string, password string) (User, error) {
	if email == "" {
		return User{}, errors.New("invalid email")
	}

	if err := db.ensureDB(); err != nil {
		return User{}, err
	}

	datbase, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, usr := range datbase.Users {
		if usr.Email == email {
			return User{}, errors.New("user with this email already exists")
		}
	}

	index := len(datbase.Users) + 1
	return User{
		Id:       index,
		Password: password,
		Email:    email,
	}, nil
}
