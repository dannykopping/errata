package errata

import (
	"github.com/dannykopping/errata/pkg/errors"
)

type DataSource interface {
	List() []errors.Error
	FindByCode(code string) errors.Error
}
