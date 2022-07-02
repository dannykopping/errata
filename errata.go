// Package errata is auto-generated by errata
// Errata Schema Version: 0.1
// Hash: e93543d95230ddb6978a8726b127c52b
package errata

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type erratum struct {
	code       string
	message    string
	categories []string
	args       map[string]interface{}
	labels     map[string]string
	guide      string

	file string
	line int

	uuid    string
	wrapped error
}

// TODO: add documentation to all public methods

// Erratum is the public interface which indicates that a given error is an Erratum.
type Erratum interface {
	// behave like a regular error
	error
	Unwrap() error

	Code() string
	Message() string
	Categories() []string
	Args() map[string]interface{}
	Guide() string
	Labels() map[string]string

	UUID() string
	HelpURL() string
}

func (e *erratum) Unwrap() error {
	return e.wrapped
}

func (e *erratum) UUID() string {
	if e.uuid == "" {
		e.uuid = generateReference(e.code)
	}
	return e.uuid
}

// Format controls the verbosity of the printed error.
func (e *erratum) Format(f fmt.State, verb rune) {
	if verb == 'v' && f.Flag('+') {
		f.Write([]byte(fmt.Sprintf("%s. For more details, see %s", e.Error(), e.HelpURL())))
		if unwrapped := e.Unwrap(); unwrapped != nil {
			if e, ok := unwrapped.(fmt.Formatter); ok {
				f.Write([]byte("\n↳ "))
				e.Format(f, verb)
			}
		}
	} else {
		f.Write([]byte(e.Error()))
	}
}

func (e *erratum) Error() string {
	return fmt.Sprintf("[errata-%s] [%s:%v] %s", e.code, e.file, e.line, e.message)
}

func (e *erratum) HelpURL() string {
	return fmt.Sprintf("https://dannykopping.github.io/errata/errata/%s", e.code)
}

func (e *erratum) Code() string {
	return e.code
}

func (e *erratum) Message() string {
	return e.message
}

func (e *erratum) Categories() []string {
	return e.categories
}

func (e *erratum) Args() map[string]interface{} {
	return e.args
}

func (e *erratum) Labels() map[string]string {
	return e.labels
}

func (e *erratum) Guide() string {
	return e.guide
}

func (e *erratum) File() string {
	return e.file
}

func (e *erratum) Line() int {
	return e.line
}

const (
	ArgumentLabelNameClashErrCode string = "argument-label-name-clash"
	CodeGenErrCode                string = "code-gen"
	FileNotFoundErrCode           string = "file-not-found"
	FileNotReadableErrCode        string = "file-not-readable"
	InvalidDatasourceErrCode      string = "invalid-datasource"
	InvalidDefinitionsErrCode     string = "invalid-definitions"
	InvalidSyntaxErrCode          string = "invalid-syntax"
	MarkdownRenderingErrCode      string = "markdown-rendering"
	ServeMethodNotAllowedErrCode  string = "serve-method-not-allowed"
	ServeSearchIndexErrCode       string = "serve-search-index"
	ServeSearchMissingTermErrCode string = "serve-search-missing-term"
	ServeUnknownCodeErrCode       string = "serve-unknown-code"
	ServeUnknownRouteErrCode      string = "serve-unknown-route"
	ServeWebUiErrCode             string = "serve-web-ui"
	TemplateExecutionErrCode      string = "template-execution"
)

type ArgumentLabelNameClashErr struct {
	erratum
}
type CodeGenErr struct {
	erratum
}
type FileNotFoundErr struct {
	erratum
}
type FileNotReadableErr struct {
	erratum
}
type InvalidDatasourceErr struct {
	erratum
}
type InvalidDefinitionsErr struct {
	erratum
}
type InvalidSyntaxErr struct {
	erratum
}
type MarkdownRenderingErr struct {
	erratum
}
type ServeMethodNotAllowedErr struct {
	erratum
}
type ServeSearchIndexErr struct {
	erratum
}
type ServeSearchMissingTermErr struct {
	erratum
}
type ServeUnknownCodeErr struct {
	erratum
}
type ServeUnknownRouteErr struct {
	erratum
}
type ServeWebUiErr struct {
	erratum
}
type TemplateExecutionErr struct {
	erratum
}

func NewArgumentLabelNameClashErr(wrapped error, key string) *ArgumentLabelNameClashErr {
	err := erratum{
		code:       ArgumentLabelNameClashErrCode,
		message:    `An error definition contains a label with the same name as an argument`,
		categories: []string{"datasource", "validation"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: `Error definitions must have labels with keys that are unique across the list of arguments`,

		args: map[string]interface{}{

			"key": key,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ArgumentLabelNameClashErr{err}
}

// GetKey returns the "key" argument for a ArgumentLabelNameClashErr instance.
func (e *ArgumentLabelNameClashErr) GetKey() interface{} {
	return e.args["key"]
}

// GetSeverity returns the "severity" label for a ArgumentLabelNameClashErr instance.
func (e *ArgumentLabelNameClashErr) GetSeverity() string {
	return "fatal"
}

func NewCodeGenErr(wrapped error) *CodeGenErr {
	err := erratum{
		code:       CodeGenErrCode,
		message:    `Code generation failed`,
		categories: []string{"codegen"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: `The provided template may contain errors`,

		args:    map[string]interface{}{},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &CodeGenErr{err}
}

// GetSeverity returns the "severity" label for a CodeGenErr instance.
func (e *CodeGenErr) GetSeverity() string {
	return "fatal"
}

func NewFileNotFoundErr(wrapped error, path string) *FileNotFoundErr {
	err := erratum{
		code:       FileNotFoundErrCode,
		message:    `File path %q is incorrect or inaccessible`,
		categories: []string{"file"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: `Ensure the given file exists and can be accessed by errata`,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &FileNotFoundErr{err}
}

// GetPath returns the "path" argument for a FileNotFoundErr instance.
func (e *FileNotFoundErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a FileNotFoundErr instance.
func (e *FileNotFoundErr) GetSeverity() string {
	return "fatal"
}

func NewFileNotReadableErr(wrapped error, path string) *FileNotReadableErr {
	err := erratum{
		code:       FileNotReadableErrCode,
		message:    `File %q is unreadable`,
		categories: []string{"file"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: `Ensure the given file has the correct permissions`,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &FileNotReadableErr{err}
}

// GetPath returns the "path" argument for a FileNotReadableErr instance.
func (e *FileNotReadableErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a FileNotReadableErr instance.
func (e *FileNotReadableErr) GetSeverity() string {
	return "fatal"
}

func NewInvalidDatasourceErr(wrapped error, path string) *InvalidDatasourceErr {
	err := erratum{
		code:       InvalidDatasourceErrCode,
		message:    `Datasource file %q is invalid`,
		categories: []string{"datasource"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: `Check the given datasource file for errors`,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &InvalidDatasourceErr{err}
}

// GetPath returns the "path" argument for a InvalidDatasourceErr instance.
func (e *InvalidDatasourceErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a InvalidDatasourceErr instance.
func (e *InvalidDatasourceErr) GetSeverity() string {
	return "fatal"
}

func NewInvalidDefinitionsErr(wrapped error, path string) *InvalidDefinitionsErr {
	err := erratum{
		code:       InvalidDefinitionsErrCode,
		message:    `One or more definitions declared in %q are invalid`,
		categories: []string{"definitions", "validation"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: ``,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &InvalidDefinitionsErr{err}
}

// GetPath returns the "path" argument for a InvalidDefinitionsErr instance.
func (e *InvalidDefinitionsErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a InvalidDefinitionsErr instance.
func (e *InvalidDefinitionsErr) GetSeverity() string {
	return "fatal"
}

func NewInvalidSyntaxErr(wrapped error, path string) *InvalidSyntaxErr {
	err := erratum{
		code:       InvalidSyntaxErrCode,
		message:    `File %q has syntax errors`,
		categories: []string{"parsing"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: ``,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &InvalidSyntaxErr{err}
}

// GetPath returns the "path" argument for a InvalidSyntaxErr instance.
func (e *InvalidSyntaxErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a InvalidSyntaxErr instance.
func (e *InvalidSyntaxErr) GetSeverity() string {
	return "fatal"
}

func NewMarkdownRenderingErr(wrapped error) *MarkdownRenderingErr {
	err := erratum{
		code:       MarkdownRenderingErrCode,
		message:    `Markdown rendering failed`,
		categories: []string{"web-ui"},
		labels: map[string]string{
			"severity": "warning",
		},
		guide: ``,

		args:    map[string]interface{}{},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &MarkdownRenderingErr{err}
}

// GetSeverity returns the "severity" label for a MarkdownRenderingErr instance.
func (e *MarkdownRenderingErr) GetSeverity() string {
	return "warning"
}

func NewServeMethodNotAllowedErr(wrapped error, method string, route string) *ServeMethodNotAllowedErr {
	err := erratum{
		code:       ServeMethodNotAllowedErrCode,
		message:    `Given HTTP method %q for requested route %q is not allowed`,
		categories: []string{"serve", "web-ui"},
		labels: map[string]string{
			"severity": "warning",
		},
		guide: ``,

		args: map[string]interface{}{

			"method": method,

			"route": route,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeMethodNotAllowedErr{err}
}

// GetMethod returns the "method" argument for a ServeMethodNotAllowedErr instance.
func (e *ServeMethodNotAllowedErr) GetMethod() interface{} {
	return e.args["method"]
}

// GetRoute returns the "route" argument for a ServeMethodNotAllowedErr instance.
func (e *ServeMethodNotAllowedErr) GetRoute() interface{} {
	return e.args["route"]
}

// GetSeverity returns the "severity" label for a ServeMethodNotAllowedErr instance.
func (e *ServeMethodNotAllowedErr) GetSeverity() string {
	return "warning"
}

func NewServeSearchIndexErr(wrapped error) *ServeSearchIndexErr {
	err := erratum{
		code:       ServeSearchIndexErrCode,
		message:    `Failed to build search index`,
		categories: []string{"serve", "web-ui", "search"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: ``,

		args:    map[string]interface{}{},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeSearchIndexErr{err}
}

// GetSeverity returns the "severity" label for a ServeSearchIndexErr instance.
func (e *ServeSearchIndexErr) GetSeverity() string {
	return "fatal"
}

func NewServeSearchMissingTermErr(wrapped error) *ServeSearchMissingTermErr {
	err := erratum{
		code:       ServeSearchMissingTermErrCode,
		message:    `Search request is missing a "term" query string parameter`,
		categories: []string{"serve", "web-ui", "search"},
		labels: map[string]string{
			"severity": "warning",
		},
		guide: ``,

		args:    map[string]interface{}{},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeSearchMissingTermErr{err}
}

// GetSeverity returns the "severity" label for a ServeSearchMissingTermErr instance.
func (e *ServeSearchMissingTermErr) GetSeverity() string {
	return "warning"
}

func NewServeUnknownCodeErr(wrapped error, code string) *ServeUnknownCodeErr {
	err := erratum{
		code:       ServeUnknownCodeErrCode,
		message:    `Cannot find error definition for given code %q`,
		categories: []string{"serve", "web-ui"},
		labels: map[string]string{
			"http_status_code": "404",
			"severity":         "warning",
		},
		guide: ``,

		args: map[string]interface{}{

			"code": code,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeUnknownCodeErr{err}
}

// GetCode returns the "code" argument for a ServeUnknownCodeErr instance.
func (e *ServeUnknownCodeErr) GetCode() interface{} {
	return e.args["code"]
}

// GetHttpStatusCode returns the "http_status_code" label for a ServeUnknownCodeErr instance.
func (e *ServeUnknownCodeErr) GetHttpStatusCode() string {
	return "404"
}

// GetSeverity returns the "severity" label for a ServeUnknownCodeErr instance.
func (e *ServeUnknownCodeErr) GetSeverity() string {
	return "warning"
}

func NewServeUnknownRouteErr(wrapped error, route string) *ServeUnknownRouteErr {
	err := erratum{
		code:       ServeUnknownRouteErrCode,
		message:    `Requested route %q not defined`,
		categories: []string{"serve", "web-ui"},
		labels: map[string]string{
			"severity": "warning",
		},
		guide: ``,

		args: map[string]interface{}{

			"route": route,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeUnknownRouteErr{err}
}

// GetRoute returns the "route" argument for a ServeUnknownRouteErr instance.
func (e *ServeUnknownRouteErr) GetRoute() interface{} {
	return e.args["route"]
}

// GetSeverity returns the "severity" label for a ServeUnknownRouteErr instance.
func (e *ServeUnknownRouteErr) GetSeverity() string {
	return "warning"
}

func NewServeWebUiErr(wrapped error, path string) *ServeWebUiErr {
	err := erratum{
		code:       ServeWebUiErrCode,
		message:    `Cannot serve web UI for datasource %q`,
		categories: []string{"serve", "web-ui"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: ``,

		args: map[string]interface{}{

			"path": path,
		},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &ServeWebUiErr{err}
}

// GetPath returns the "path" argument for a ServeWebUiErr instance.
func (e *ServeWebUiErr) GetPath() interface{} {
	return e.args["path"]
}

// GetSeverity returns the "severity" label for a ServeWebUiErr instance.
func (e *ServeWebUiErr) GetSeverity() string {
	return "fatal"
}

func NewTemplateExecutionErr(wrapped error) *TemplateExecutionErr {
	err := erratum{
		code:       TemplateExecutionErrCode,
		message:    `Error in template execution`,
		categories: []string{"codegen"},
		labels: map[string]string{
			"severity": "fatal",
		},
		guide: ``,

		args:    map[string]interface{}{},
		wrapped: wrapped,
	}

	addCaller(&err)
	return &TemplateExecutionErr{err}
}

// GetSeverity returns the "severity" label for a TemplateExecutionErr instance.
func (e *TemplateExecutionErr) GetSeverity() string {
	return "fatal"
}
func addCaller(err *erratum) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		paths := strings.Split(file, string(os.PathSeparator))
		segments := 2
		if len(paths) < segments {
			segments = 1
		}
		err.file = filepath.Join(paths[len(paths)-segments:]...)
		err.line = line
	}
}

func generateReference(code string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(code+time.Now().Format(time.RFC3339Nano))))
}
