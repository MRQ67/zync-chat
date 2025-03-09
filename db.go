package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./zyncchat.db")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
	return db
}

func RegisterUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	return err
}

func AuthenticateUser(db *sql.DB, username, password string) error {
	var storedHash string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}
