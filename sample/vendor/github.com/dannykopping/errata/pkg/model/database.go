package model

import "fmt"

type Database struct {
	Version string `yaml:"version"`

	Categories interface{} `yaml:"categories"`

	Errors []*Error `yaml:"errors"`
}

type Error struct {
	Code       string
	Message    string   `yaml:"message"`
	Cause      string   `yaml:"cause,omitempty"`
	Categories []string `yaml:"categories,omitempty,flow"`

	External *Error `yaml:"external,omitempty"`
	HTTP     *HTTP  `yaml:"http,omitempty"`
}

type HTTP struct {
	Code int `yaml:"code"`
}

var (
	DatabaseUninitialized = Error{
		Code:    "db_uninitialized",
		Message: "Errata database is not initialized",
	}
	UnknownErrataCode = Error{
		Code:    "unknown_errata_code",
		Message: "Could not locate given errata code",
	}
)

func (db *Database) FindByCode(code string) error {
	if db == nil {
		return &DatabaseUninitialized
	}

	for _, e := range db.Errors {
		if e.Code == code {
			return e
		}
	}

	fmt.Printf("could not find code: %q\n", code)
	return &UnknownErrataCode
}

func (e *Error) Error() string {
	return fmt.Sprintf(">>>%s", e.Code)
}