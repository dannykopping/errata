package errata

import (
	"fmt"

	"github.com/dannykopping/errata/pkg/errors"
)

var (
	source DataSource
)

var (
	InvalidDataSource = map[string]interface{}{
		"Code": "invalid_datasource",
	}
)

func RegisterDataSource(ds DataSource) error {
	if ds == nil {
		// TODO: use errata
		return fmt.Errorf("invalid datasource")
	}

	source = ds
	return nil
}

func New(code string) errors.ErrorDefinition {
	return source.FindByCode(code)
}
