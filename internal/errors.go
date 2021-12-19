package errors

import "fmt"

type Error struct {
	Code       string            `yaml:"code"`
	Message    string            `yaml:"message"`
	Cause      string            `yaml:"cause,omitempty"`
	Categories []string          `yaml:"categories,omitempty,flow"`
	Labels     map[string]string `yaml:"labels,omitempty,flow"`
	Args       []interface{}     `yaml:"args"`

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
		return fmt.Sprintf("[%s] %s", e.String(), w)
	}

	return e.String()
}

func (e Error) String() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

// FileNotFound: YML file path is incorrect or inaccessible
func FileNotFound(path interface{}) Error {
	return Error{
		Code:       "file_not_found",
		Message:    "YML file path is incorrect or inaccessible",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Args:       []interface{}{path},
	}
}

// FileNotReadable: YML file path is unreadable
func FileNotReadable(path interface{}) Error {
	return Error{
		Code:       "file_not_readable",
		Message:    "YML file path is unreadable",
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Args:       []interface{}{path},
	}
}

// SyntaxError: YML file is malformed
func SyntaxError() Error {
	return Error{
		Code:       "syntax_error",
		Message:    "YML file is malformed",
		Cause:      "",
		Categories: []string{"parsing"},
		Labels:     map[string]string{},
		Args:       []interface{}{},
	}
}

// CodeGenError: Code generation failed
func CodeGenError() Error {
	return Error{
		Code:       "code_gen_error",
		Message:    "Code generation failed",
		Cause:      "",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		Args:       []interface{}{},
	}
}
