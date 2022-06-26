package store

import (
	"crypto/md5"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	Email        string
	PasswordHash string
	Spam         bool
	Abuse        bool
}

func NewUsersStore(path string) (*UserStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		// TODO wrap error
		return nil, err
	}

	return &UserStore{db: db}, nil
}

func (d *UserStore) GetUser(email string, password string) (*User, error) {
	var u User

	r := d.db.QueryRow("SELECT * FROM users WHERE email = :email AND password_hash = :hash",
		sql.Named("email", email),
		sql.Named("hash", d.hash(password)),
	)
	err := r.Err()
	if err != nil {
		// TODO wrap error
		return nil, err
	}

	err = r.Scan(&u.Email, &u.PasswordHash, &u.Spam, &u.Abuse)
	if err != nil {
		if err == sql.ErrNoRows {
		}
		// TODO wrap error
		return nil, err
	}

	return &u, nil
}

func (d *UserStore) hash(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
