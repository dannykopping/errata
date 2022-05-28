package errata

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	I18nDir string
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

	// TODO: error handling
	i18nMap, err := buildI18nMap(data)
	if err != nil {
		panic(err)
	}

	tmplData := pongo2.Context{
		"Package": data.Package,
		"Errors":  source.List(),
		"I18n":    i18nMap,
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

func buildI18nMap(data CodeGen) (map[string]map[string]ErrorDefinition, error) {
	// TODO improve error handling here

	m := make(map[string]map[string]ErrorDefinition)

	_, err := os.Stat(data.I18nDir)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewFileNotFoundErr(err)
		}

		return nil, NewInvalidDatasourceErr(err)
	}

	err = filepath.WalkDir(data.I18nDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		source, err := NewFileDatasource(path)
		if err != nil {
			return err
		}

		locale := strings.Split(d.Name(), ".")[0]

		m[locale] = make(map[string]ErrorDefinition, len(source.List()))
		for code, e := range source.List() {
			m[locale][code] = e
		}

		return nil
	})

	if err != nil {
		return nil, NewInvalidDatasourceErr(err)
	}

	return m, nil
}
