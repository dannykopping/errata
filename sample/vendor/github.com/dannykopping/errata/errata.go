package errata

import "github.com/dannykopping/errata/pkg/model"

type DataSource interface {
	FindByCode(string) error
}

type Metadata struct {
	Key string
	Value interface{}
}

var source DataSource

var (
	InvalidDataSource = model.Error{
		Code:    "invalid_datasource",
		Message: "Given errata datasource is invalid",
	}
)

func RegisterSource(ds DataSource) error {
	if ds == nil {
		return &InvalidDataSource
	}

	source = ds
	return nil
}

func New(code string) error {
	return source.FindByCode(code)
}
