package errata

import (
	"io"
)

func Generate(data CodeGenConfig, w io.Writer) error {
	source, err := NewHCLDatasource(data.Source)
	if err != nil {
		return err
	}

	if err := source.Validate(); err != nil {
		return NewInvalidDefinitionsErr(err, data.Source)
	}

	loader, err := loaderFromPath(data.Template)
	if err != nil {
		return err
	}

	tmplData := map[string]interface{}{
		"Package": data.Package,
		"Options": source.Options(),
		"Errors":  source.List(),
		"Version": source.Version(),
	}

	renderer := createPongo2Renderer(loader)
	if err = renderer.render(tmplData, w); err != nil {
		return NewCodeGenErr(err)
	}

	return nil
}
