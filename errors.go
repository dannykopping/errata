// NOTE: this file is auto-generated by errata
package errata

import "fmt"

type Error struct {
	Code       string
	Message    string
	Cause      string
	Categories []string
	Labels     map[string]string
	Interfaces interfaces

	inner error
}

type interfaces struct {
	HTTPResponseCode int
	ShellExitCode    int
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

func CodeGenError() Error {
	return Error{
		Code:       "code_gen_error",
		Message:    "Code generation failed",
		Cause:      "",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func FileNotFound(path interface{}) Error {
	return Error{
		Code:       "file_not_found",
		Message:    fmt.Sprintf("YML file [%s] is incorrect or inaccessible", path),
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func FileNotReadable(path interface{}) Error {
	return Error{
		Code:       "file_not_readable",
		Message:    fmt.Sprintf("YML file [%s] is unreadable", path),
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func InvalidDatasource(datasource interface{}) Error {
	return Error{
		Code:       "invalid_datasource",
		Message:    fmt.Sprintf("Configured datasource [%s] is invalid", datasource),
		Cause:      "",
		Categories: []string{"init"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func SyntaxError(path interface{}) Error {
	return Error{
		Code:       "syntax_error",
		Message:    fmt.Sprintf("YML file [%s] is malformed", path),
		Cause:      "",
		Categories: []string{"parsing"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func TemplateExecution() Error {
	return Error{
		Code:       "template_execution",
		Message:    "Error in template execution",
		Cause:      "Possible use of missing or renamed field",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		Interfaces: interfaces{HTTPResponseCode: 3, ShellExitCode: 3},
	}
}

func TemplateNotFound(path interface{}) Error {
	return Error{
		Code:       "template_not_found",
		Message:    fmt.Sprintf("Template path [%s] is incorrect or inaccessible", path),
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func TemplateNotReadable(path interface{}) Error {
	return Error{
		Code:       "template_not_readable",
		Message:    fmt.Sprintf("Template path [%s] is unreadable", path),
		Cause:      "",
		Categories: []string{"file"},
		Labels:     map[string]string{},
		Interfaces: interfaces{},
	}
}

func TemplateSyntax() Error {
	return Error{
		Code:       "template_syntax",
		Message:    "Syntax error in template",
		Cause:      "",
		Categories: []string{"codegen"},
		Labels:     map[string]string{},
		Interfaces: interfaces{ShellExitCode: 2},
	}
}