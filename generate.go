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
	//source, err := NewFileDatasource(data.File)
	source, err := NewHCLDatasource(data.File)
	if err != nil {
		return err
	}

	if err := source.Validate(); err != nil {
		return NewInvalidDefinitionsErr(err, data.File)
	}

	// TODO: support built-in and external templates
	//		-> built-in: -template=golang
	//		-> external: -template=my-template.tmpl
	file := fmt.Sprintf("%s.tmpl", data.Lang)
	path := fmt.Sprintf("templates/%s", file)

	tmplData := pongo2.Context{
		"Package": data.Package,
		"Options": source.Options(),
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

	templateSet := pongo2.NewSet("blah", pongo2.NewFSLoader(templates))
	pongo2.SetAutoescape(false)

	b, err := templates.ReadFile(path)
	if err != nil {
		return NewTemplateNotReadableErr(err)
	}

	tmpl, err := templateSet.FromBytes(b)
	if err != nil {
		return NewTemplateSyntaxErr(err)
	}

	if err := tmpl.ExecuteWriter(tmplData, w); err != nil {
		return NewTemplateExecutionErr(err)
	}

	return nil
}
