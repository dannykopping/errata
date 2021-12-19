package main

import "fmt"

type Error struct {
	Code       string
	Message    string
	Cause      string
	Categories []string
	Labels     map[string]string

	inner error
}

func (e Error) Wrap(err error) error {
	if err == nil {
		return nil
	}

	e.inner = err
	return e
}

func (e Error) Unwrap() error {
	return e.inner
}

func (e Error) Error() string {
	if w := e.Unwrap(); w != nil {
		return fmt.Sprintf("[code: %s] %s", e.Code, w)
	}

	return e.Code
}

func TemplateExecution() Error {
	return Error{
		Code:       "template_execution",
		Message:    "Error in template execution",
		Cause:      "Possible use of missing or renamed field",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
	}
}

func TemplateNotFound(path interface{}) Error {
	return Error{
		Code:       "template_not_found",
		Message:    fmt.Sprintf("Template path [%s] is incorrect or inaccessible", path),
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
	}
}

func TemplateNotReadable() Error {
	return Error{
		Code:       "template_not_readable",
		Message:    "Template path is unreadable",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
	}
}

func TemplateSyntax() Error {
	return Error{
		Code:       "template_syntax",
		Message:    "Syntax error in template",
		Cause:      "",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
	}
}
