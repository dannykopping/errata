package errata

import (
	"io"
	"os"

	internal "github.com/dannykopping/errata/internal"
	"github.com/dannykopping/errata/pkg/errors"
	"gopkg.in/yaml.v2"
)

type fileDatasource struct {
	Version string `yaml:"version"`

	Errors []errors.Error `yaml:"errors"`
}

func NewFileDatasource(path string) (DataSource, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, internal.FileNotFound().Wrap(err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, internal.FileNotReadable().Wrap(err)
	}

	db, err := parse(f)
	if err != nil {
		return nil, internal.SyntaxError().Wrap(err)
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

func (db *fileDatasource) FindByCode(code string) errors.Error {
	for _, e := range db.Errors {
		if e.Code == code {
			return e
		}
	}

	// if we cannot find the error by code, create one
	return errors.Error{
		Code: code,
	}
}

func (db *fileDatasource) List() []errors.Error {
	return db.Errors
}
