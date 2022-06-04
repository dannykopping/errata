package errata

import (
	"embed"
	"fmt"
	"io"

	"github.com/flosch/pongo2/v5"
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

func Generate(data CodeGen, w io.Writer) error {
	source, err := NewFileDatasource(data.File)
	if err != nil {
		return err
	}

	// TODO: support built-in and external templates
	//		-> built-in: -template=golang
	//		-> external: -template=my-template.tmpl
	file := fmt.Sprintf("%s.tmpl", data.Lang)
	path := fmt.Sprintf("templates/%s/%s", data.Lang, file)

	tmplData := pongo2.Context{
		"Package": data.Package,
		"Errors":  source.List(),
		"Version": source.Version(),
	}

	_, err = templates.Open(path)
	if err != nil {
		return NewTemplateNotFoundErr(err)
	}

	pongo2.RegisterFilter("constantize", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsValue(strcase.ToCamel(in.String())), nil
	})

	pongo2.SetAutoescape(false)

	tmpl, err := pongo2.FromFile(path)
	if err != nil {
		return NewTemplateSyntaxErr(err)
	}

	if err := tmpl.ExecuteWriter(tmplData, w); err != nil {
		return NewTemplateExecutionErr(err)
	}

	return nil
}
