package errata

import (
	"io"

	"gopkg.in/yaml.v2"
)

var db *Database

func Parse(reader io.Reader) (*Database, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
