package errata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var input = `
---
version: '0.1'

errors:
  template_not_found:
    message: Template path [%s] is incorrect or inaccessible
    categories: [ file ]
    args:
      - path
  template_syntax:
    message: Syntax error in template
    categories: [ codegen ]
`

func TestUnmarshal(t *testing.T) {
	var db fileDatasource
	if err := yaml.Unmarshal([]byte(input), &db); err != nil {
		fmt.Println(err.Error())
	}

	require.Len(t, db.Definitions, 2)

	code := "template_not_found"
	def := db.Definitions[code]

	require.NotEmpty(t, def)
	require.Equal(t, code, def.Code)
	require.Equal(t, "Template path [%s] is incorrect or inaccessible", def.Definition["message"])
	require.Equal(t, []interface{}{"file"}, def.Definition["categories"])
}
