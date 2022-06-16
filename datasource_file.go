package errata

//
//import (
//	"crypto/md5"
//	"fmt"
//	"io"
//	"os"
//	"sort"
//	"time"
//
//	"gopkg.in/yaml.v2"
//)
//
//type fileDatasource struct {
//	source []byte
//
//	SchemaVersion string
//	Definitions   map[string]ErrorDefinition
//	Opts          ErrorOptions
//}
//
//func (e *fileDatasource) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	var s struct {
//		Version string                            `yaml:"version"`
//		Options ErrorOptions                      `yaml:"options"`
//		Errors  map[string]map[string]interface{} `yaml:"errors"`
//	}
//
//	err := unmarshal(&s)
//	if err != nil {
//		return err
//	}
//
//	e.SchemaVersion = fmt.Sprintf("%v", s.Version)
//	e.Opts = s.Options
//	e.Definitions = make(map[string]ErrorDefinition, len(s.Errors))
//
//	// sort map keys so generated code can be idempotent
//	var codes []string
//	for code := range s.Errors {
//		codes = append(codes, code)
//	}
//
//	sort.Strings(codes)
//
//	for _, code := range codes {
//		e.Definitions[code] = ErrorDefinition{
//			Code:       code,
//			Definition: s.Errors[code],
//		}
//	}
//
//	return nil
//}
//
//func (e *fileDatasource) Version() string {
//	return fmt.Sprintf("%x/%s", md5.Sum(e.source), time.Now().Format(time.RFC3339))
//}
//
//func NewFileDatasource(path string) (DataSource, error) {
//	if _, err := os.Stat(path); err != nil {
//		return nil, NewFileNotFoundErr(err)
//	}
//
//	f, err := os.Open(path)
//	if err != nil {
//		return nil, NewFileNotReadableErr(err)
//	}
//
//	db, err := parse(f)
//	if err != nil {
//		return nil, NewInvalidSyntaxErr(err)
//	}
//
//	return db, nil
//}
//
//func parse(reader io.Reader) (*fileDatasource, error) {
//	b, err := io.ReadAll(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	var db *fileDatasource
//
//	err = yaml.Unmarshal(b, &db)
//	if err != nil {
//		return nil, err
//	}
//
//	db.source = b
//	return db, nil
//}
//
//func (db *fileDatasource) FindByCode(code string) ErrorDefinition {
//	for _, e := range db.Definitions {
//		if e.Code == code {
//			return e
//		}
//	}
//
//	// if we cannot find the error by code, create one
//	return ErrorDefinition{
//		Code: code,
//	}
//}
//
//func (db *fileDatasource) List() map[string]ErrorDefinition {
//	return db.Definitions
//}
//
//func (db *fileDatasource) Options() ErrorOptions {
//	return db.Opts
//}
