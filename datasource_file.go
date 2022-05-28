package errata

import (
	"fmt"
	"io"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type fileDatasource struct {
	Version     string
	Definitions map[string]ErrorDefinition
}

func (e *fileDatasource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s struct {
		Version string                            `yaml:"version"`
		Errors  map[string]map[string]interface{} `yaml:"errors"`
	}

	err := unmarshal(&s)
	if err != nil {
		return err
	}

	e.Version = fmt.Sprintf("%v", s.Version)
	e.Definitions = make(map[string]ErrorDefinition, len(s.Errors))

	// sort map keys so generated code can be idempotent
	var codes []string
	for code := range s.Errors {
		codes = append(codes, code)
	}

	sort.Strings(codes)

	for _, code := range codes {
		e.Definitions[code] = ErrorDefinition{
			Code:       code,
			Definition: s.Errors[code],
		}
	}

	return nil
}

func NewFileDatasource(path string) (DataSource, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, NewFileNotFoundErr(err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, NewFileNotReadableErr(err)
	}

	db, err := parse(f)
	if err != nil {
		return nil, NewInvalidSyntaxErr(err)
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

func (db *fileDatasource) FindByCode(code string) ErrorDefinition {
	for _, e := range db.Definitions {
		if e.Code == code {
			return e
		}
	}

	// if we cannot find the error by code, create one
	return ErrorDefinition{
		Code: code,
	}
}

func (db *fileDatasource) List() map[string]ErrorDefinition {
	return db.Definitions
}
