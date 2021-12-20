package errata

import (
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	//go:embed templates/*
	templates embed.FS
)

type Tmpl struct {
	Package string
	Errors  map[string]ErrorDefinition
}

type CodeGen struct {
	File    string
	Lang    string
	Package string
}

func Generate(data CodeGen) error {
	source, err := NewFileDatasource(data.File)
	if err != nil {
		return err
	}

	// TODO: support built-in and external templates
	//		-> built-in: -template=golang
	//		-> external: -template=my-template.tmpl
	file := fmt.Sprintf("%s.tmpl", data.Lang)
	path := fmt.Sprintf("templates/%s/%s", data.Lang, file)

	tmplData := Tmpl{
		Package: data.Package,
		Errors:  source.List(),
	}

	_, err = templates.Open(path)
	if err != nil {
		return TemplateNotFound(path).Wrap(err)
	}

	tmpl, err := template.New(file).
		Funcs(template.FuncMap{
			"constantize": strcase.ToCamel,
			"hasNext": func(index int, len int) bool {
				return index+1 < len
			},
			// credit: https://stackoverflow.com/a/18276968
			"input": func(in ...interface{}) (map[string]interface{}, error) {
				if len(in)%2 != 0 {
					return nil, fmt.Errorf("uneven set of parameter pairs")
				}

				var out = make(map[string]interface{}, len(in)/2)
				for i := 0; i < len(in); i += 2 {
					key, ok := in[i].(string)
					if !ok {
						return nil, fmt.Errorf("keys must be strings")
					}
					out[key] = in[i+1]
				}

				return out, nil
			},
			"debug": func(val interface{}) string {
				return ""
			},
		}).
		ParseFS(templates, "templates/**/*")

	if err != nil {
		return TemplateSyntax().Wrap(err)
	}

	return TemplateExecution().Wrap(
		tmpl.Execute(os.Stdout, tmplData),
	)
}
