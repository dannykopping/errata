package errors

import "github.com/dannykopping/errata/pkg/errors"

func FileNotFound() errors.Error {
	return errors.Error{
		Code:       "file_not_found",
		Message:    "YML file path is incorrect or inaccessible",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 1},
	}
}

func FileNotReadable() errors.Error {
	return errors.Error{
		Code:       "file_not_readable",
		Message:    "YML file path is unreadable",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 1},
	}
}

func SyntaxError() errors.Error {
	return errors.Error{
		Code:       "syntax_error",
		Message:    "YML file is malformed",
		Cause:      "",
		Categories: []string{"parsing"},
		Labels:     map[string]string{},
		HTTP:       errors.HTTP{StatusCode: 0},
		Shell:      errors.Shell{ExitCode: 2},
	}
}
