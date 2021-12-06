package main

import "github.com/dannykopping/errata/pkg/errors"

func TemplateNotFoundErr() errors.Error {
	return errors.Error{
		Code:       "template_not_found",
		Message:    "Template path is incorrect or inaccessible",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 1},
	}
}

func TemplateNotReadableErr() errors.Error {
	return errors.Error{
		Code:       "template_not_readable",
		Message:    "Template path is unreadable",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 1},
	}
}

func TemplateSyntaxErr() errors.Error {
	return errors.Error{
		Code:       "template_syntax",
		Message:    "Syntax error in template",
		Cause:      "",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 2},
	}
}

func TemplateExecutionErr() errors.Error {
	return errors.Error{
		Code:       "template_execution",
		Message:    "Error in template execution",
		Cause:      "Possible use of missing or renamed field",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 3},
	}
}
