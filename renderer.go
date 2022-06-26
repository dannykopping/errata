package errata

import (
	"io"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/iancoleman/strcase"
)

type pongo2Renderer struct {
	loader *templateLoader
}

func preparePongo2() {
	pongo2.RegisterFilter("constantize", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		return pongo2.AsValue(strcase.ToCamel(in.String())), nil
	})

	// backticks can't be escaped in golang multi-line strings
	pongo2.RegisterFilter("escape_backtick", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		splits := strings.Split(in.String(), "`")
		return pongo2.AsValue(strings.Join(splits, "` + \"`\" + `")), nil
	})
	pongo2.SetAutoescape(false)
}

func createPongo2Renderer(loader *templateLoader) *pongo2Renderer {
	preparePongo2()
	return &pongo2Renderer{
		loader: loader,
	}
}

func (p *pongo2Renderer) getTemplate() (*pongo2.Template, error) {
	if !p.loader.builtin {
		tmpl, err := pongo2.DefaultSet.FromFile(p.loader.path)
		if err != nil {
			return nil, NewInvalidSyntaxErr(err, p.loader.path)
		}

		return tmpl, nil
	}

	b, err := templates.ReadFile(p.loader.path)
	if err != nil {
		return nil, NewFileNotReadableErr(err, p.loader.path)
	}

	tmpl, err := embeddedFS.FromBytes(b)
	if err != nil {
		return nil, NewInvalidSyntaxErr(err, p.loader.path)
	}

	return tmpl, nil
}

func (p *pongo2Renderer) render(ctx map[string]interface{}, w io.Writer) error {
	tmpl, err := p.getTemplate()
	if err != nil {
		return err
	}

	err = tmpl.ExecuteWriter(ctx, w)
	if err != nil {
		return NewTemplateExecutionErr(err)
	}

	return nil
}
