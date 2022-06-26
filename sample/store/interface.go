package store

type Store interface {
	GetUser(email string, password string) (*User, error)
}
