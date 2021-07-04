package errata

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type fileDatasource struct {
	Version string `yaml:"version"`

	Errors []*Error `yaml:"errors"`
}

var (
	fileOpenError = Error{
		Code:    "db_file_open",
		Message: "Errata database file cannot be opened",
	}
	fileParseError = Error{
		Code:    "db_file_parse",
		Message: "Errata database file cannot be parsed",
	}
)

func NewFileDatasource(path string) (*fileDatasource, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("db open error: %s\n", err)
		return nil, fileOpenError
	}

	db, err := parse(f)
	if err != nil {
		fmt.Printf("db parse error: %s\n", err)
		return nil, fileParseError
	}

	return db, nil
}

func parse(reader io.Reader) (*fileDatasource, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var db *fileDatasource

	err = yaml.Unmarshal(bytes, &db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *fileDatasource) FindByCode(code string) Error {
	if db == nil {
		return DatasourceUninitializedError
	}
	for _, e := range db.Errors {
		if e.Code == code {
			return *e
		}
	}

	// if we cannot find the error by code, create one
	return Error{
		Code: code,
	}
}
