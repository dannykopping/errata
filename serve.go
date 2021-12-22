package errata

import (
	"embed"
	"os"
	"text/template"
)

var (
	//go:embed web/*
	web embed.FS
)

func Serve(data CodeGen) error {
	source, err := NewFileDatasource(data.File)
	if err != nil {
		return err
	}

	tmplData := Tmpl{
		Package: data.Package,
		Errors:  source.List(),
	}

	tmpl, err := template.New("index.gohtml").
		ParseFS(web, "web/*")

	if err != nil {
		return TemplateSyntax().Wrap(err)
	}

	return TemplateExecution().Wrap(
		tmpl.Execute(os.Stdout, tmplData),
	)
}
