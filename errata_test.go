package errata

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Format(t *testing.T) {
	tests := []struct {
		name           string
		err            Erratum
		expectedOutput *regexp.Regexp
	}{
		{
			name:           "single error, no args",
			err:            NewCodeGenErr(nil),
			expectedOutput: regexp.MustCompile(`\[errata-code-gen] <[^>]+> Code generation failed. For more details, see .*/errata-code-gen`),
		},
		{
			name:           "single error, with single arg",
			err:            NewInvalidDefinitionsErr(nil, "foo"),
			expectedOutput: regexp.MustCompile(`\[errata-invalid-definitions] <[^>]+> One or more definitions declared in are invalid \(path="foo"\). For more details, see .*/errata-invalid-definitions`),
		},
		{
			name:           "single error, with multiple args",
			err:            NewServeMethodNotAllowedErr(nil, "foo", "GET"),
			expectedOutput: regexp.MustCompile(`\[errata-serve-method-not-allowed] <[^>]+> Given HTTP method for requested route is not allowed \(route="foo", method="GET"\). For more details, see .*/errata-serve-method-not-allowed`),
		},
		{
			name: "wrapped Erratum",
			err:  NewCodeGenErr(NewInvalidSyntaxErr(nil, "foo")),
			expectedOutput: regexp.MustCompile(`\[errata-code-gen] <[^>]+> Code generation failed. For more details, see .*/errata-code-gen
↳ \[errata-invalid-syntax] <[^>]+> File has syntax errors \(path="foo"\). For more details, see .*/errata-invalid-syntax`),
		},
		{
			name: "wrapped regular error",
			err:  NewCodeGenErr(fmt.Errorf("this is a regular error")),
			expectedOutput: regexp.MustCompile(`\[errata-code-gen] <[^>]+> Code generation failed. For more details, see .*/errata-code-gen
↳ this is a regular error`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Regexp(t, tt.expectedOutput, fmt.Sprintf("%+v", tt.err))
		})
	}
}
