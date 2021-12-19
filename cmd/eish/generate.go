package main

import (
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/pkg/errors"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed templates/*
	templates embed.FS

	edsFile string
	lang    string
	pkg     string
)

type Tmpl struct {
	Package string
	Errors  map[string]errors.ErrorDefinition
}

func generate(_ *cli.Context) error {
	source, err := errata.NewFileDatasource(edsFile)
	if err != nil {
		return err
	}

	// TODO: support built-in and external templates
	//		-> built-in: -template=golang
	//		-> external: -template=my-template.tmpl
	file := fmt.Sprintf("%s.tmpl", lang)
	path := fmt.Sprintf("templates/%s/%s", lang, file)

	data := Tmpl{
		Package: pkg,
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
		}).
		ParseFS(templates, "templates/**/*")

	if err != nil {
		return TemplateSyntax().Wrap(err)
	}

	return TemplateExecution().Wrap(
		tmpl.Execute(os.Stdout, data),
	)
}
