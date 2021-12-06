package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/dannykopping/errata"
	"github.com/dannykopping/errata/pkg/errors"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

var (
	//go:embed templates
	templates embed.FS

	edsFile string
	lang    string
	pkg     string
)

type Tmpl struct {
	Package string
	Errors  []errors.Error
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
	path := fmt.Sprintf("templates/%s", file)

	data := Tmpl{
		Package: pkg,
		Errors:  source.List(),
	}

	_, err = templates.Open(path)
	if err != nil {
		return TemplateNotFoundErr().Wrap(err)
	}

	tmpl, err := template.New(file).
		Funcs(template.FuncMap{
			"constantize": strcase.ToCamel,
			"quote": func(in []string) []string {
				var out []string
				for _, s := range in {
					out = append(out, fmt.Sprintf("%q", s))
				}
				return out
			},
			"list": func(in []string) string {
				return strings.Join(in, ", ")
			},
		}).
		ParseFS(templates, path)

	if err != nil {
		return TemplateSyntaxErr().Wrap(err)
	}

	return TemplateExecutionErr().Wrap(
		tmpl.Execute(os.Stdout, data),
	)
}
