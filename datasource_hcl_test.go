package errata

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHCLDatasource(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		list        map[string]errorDefinition
		expectedErr func(err error) (error, bool)
	}{
		{
			name:    "basic",
			fixture: "fixtures/basic.hcl",
			list: map[string]errorDefinition{
				"code-1": {
					Code:    "code-1",
					Message: "This is a basic error",
				},
			},
		},
		{
			name:    "multiple errors defined",
			fixture: "fixtures/multiple.hcl",
			list: map[string]errorDefinition{
				"code-1": {
					Code:    "code-1",
					Message: "This is a basic error",
				},
				"code-2": {
					Code:    "code-2",
					Message: "This is another basic error",
					Args: []arg{
						{
							Name: "first",
							Type: "string",
						},
						{
							Name: "second",
							Type: "bool",
						},
					},
				},
				"code-3": {
					Code:    "code-3",
					Message: "This one has a guide file",
					Guide:   "file://fixtures/guide.md",
				},
			},
		},
		{
			name:    "invalid syntax",
			fixture: "fixtures/invalid-syntax.hcl",
			expectedErr: func(err error) (error, bool) {
				var expected InvalidSyntaxErr
				ok := errors.As(err, &expected)
				return expected, ok
			},
		},
		{
			name:    "missing database",
			fixture: "fixtures/cantfindme.hcl",
			expectedErr: func(err error) (error, bool) {
				var expected FileNotFoundErr
				ok := errors.As(err, &expected)
				return expected, ok
			},
		},
		{
			name:    "empty database",
			fixture: "fixtures/empty.hcl",
			expectedErr: func(err error) (error, bool) {
				var expected InvalidDatasourceErr
				ok := errors.As(err, &expected)
				return expected, ok
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewHCLDatasource(tt.fixture)
			if err != nil {
				// yes, this would be a bit easier with generics, but I don't see this as a compelling enough
				// reason to make the lib depend on >=1.18
				expected, ok := tt.expectedErr(err)
				t.Log(err.Error())
				assert.Truef(t, ok, "Expecting error of type %T", expected)

				return
			}

			assert.Equal(t, tt.list, ds.List())
		})
	}
}

func TestHclDatasource_Options(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		opts    errorOptions
	}{
		{
			name:    "no options",
			fixture: "fixtures/no-opts.hcl",
		},
		{
			name:    "with options",
			fixture: "fixtures/with-opts.hcl",
			opts: errorOptions{
				Description: "This is a description",
				BaseURL:     "https://dannykopping.github.io/errata/",
				Prefix:      "err-",
				Imports: []string{
					"fmt",
					"github.com/hashicorp/hcl/v2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewHCLDatasource(tt.fixture)
			assert.NoError(t, err)

			assert.Equal(t, tt.opts, ds.Options())
		})
	}
}
