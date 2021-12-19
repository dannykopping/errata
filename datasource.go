package errata

import "github.com/dannykopping/errata/pkg/errors"

type DataSource interface {
	List() map[string]errors.ErrorDefinition
	FindByCode(code string) errors.ErrorDefinition
}
