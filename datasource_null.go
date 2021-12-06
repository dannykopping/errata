package errata

import (
	"github.com/dannykopping/errata/pkg/errors"
)

type NullDataSource struct{}

func NewNullDataSource() DataSource {
	return &NullDataSource{}
}

func (n *NullDataSource) FindByCode(code string) errors.Error {
	return errors.Error{
		Code: code,
	}
}

func (n *NullDataSource) List() []errors.Error {
	return []errors.Error{}
}
