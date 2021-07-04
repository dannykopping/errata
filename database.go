package errata

import "fmt"

type Database struct {
	Version string `yaml:"version"`

	Categories interface{} `yaml:"categories"`

	Errors []*Error `yaml:"errors"`
}

type HTTP struct {
	StatusCode int `yaml:"status_code"`
}

type Shell struct {
	ExitCode int `yaml:"exit_code"`
}

var (
	DatabaseFileOpen = Error{
		Code:    "db_file_open",
		Message: "Errata database file cannot be opened",
	}
	DatabaseFileParse = Error{
		Code:    "db_file_parse",
		Message: "Errata database file cannot be parsed",
	}
	DatabaseUninitialized = Error{
		Code:    "db_uninitialized",
		Message: "Errata database is not initialized",
	}
	UnknownErrataCode = Error{
		Code:    "unknown_errata_code",
		Message: "Could not locate given errata code",
	}
)

func (db *Database) FindByCode(code string) Error {
	if db == nil {
		return DatabaseUninitialized
	}

	for _, e := range db.Errors {
		if e.Code == code {
			return *e
		}
	}

	fmt.Printf("could not find code: %q\n", code)
	return UnknownErrataCode
}
