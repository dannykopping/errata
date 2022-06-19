package errata

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v5"
)

var (
	//go:embed templates/*
	templates  embed.FS
	embeddedFS = pongo2.NewSet("embedded", pongo2.NewFSLoader(templates))
)

type templateLoader struct {
	path    string
	builtin bool
}

func loaderFromPath(given string) (*templateLoader, error) {
	// first check if this is a built-in template (no file extension, language name only)
	if filepath.Ext(given) == "" {
		file := fmt.Sprintf("%s.tmpl", given)
		path := fmt.Sprintf("templates/%s", file)
		_, err := templates.Open(path)
		if err != nil {
			return nil, NewTemplateNotFoundErr(err)
		}

		return &templateLoader{
			path:    path,
			builtin: true,
		}, nil
	}

	// next try resolve the path literally
	_, err := os.Stat(given)
	if errors.Is(err, os.ErrNotExist) {
		return nil, NewTemplateNotFoundErr(err)
	}

	return &templateLoader{
		path:    given,
		builtin: false,
	}, nil
}
