package errata

import (
	"github.com/dannykopping/errata/pkg/errors"
)

var (
	source DataSource = NewNullDataSource()
)

var (
	InvalidDataSource = errors.Error{
		Code:    "invalid_datasource",
		Message: "Given Errata datasource is invalid",
	}
)

func RegisterDataSource(ds DataSource) error {
	if ds == nil {
		return InvalidDataSource
	}

	source = ds
	return nil
}

func New(code string) errors.Error {
	return source.FindByCode(code)
}
