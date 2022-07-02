package errata

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

type hclDatasource struct {
	sync.Mutex

	source []byte
	path   string
	list   map[string]errorDefinition

	Version string `hcl:"version"`
	// options are, well, optional - and HCL only allows optional blocks as pointer types
	Opts   *errorOptions        `hcl:"options,block"`
	Errors []hclErrorDefinition `hcl:"error,block"`
}

type hclErrorDefinition struct {
	Code       string            `hcl:",label"`
	Message    string            `hcl:"message"`
	Cause      string            `hcl:"cause,optional"`
	Categories []string          `hcl:"categories,optional"`
	Args       []cty.Value       `hcl:"args,optional"`
	Labels     map[string]string `hcl:"labels,optional"`
	Guide      string            `hcl:"guide,optional"`

	// TODO use Remain to allow custom fields being defined
	//Remain     hcl.Body          `hcl:",remain"`
}

type errorOptions struct {
	Prefix      string   `hcl:"prefix,optional"`
	Imports     []string `hcl:"imports,optional"`
	BaseURL     string   `hcl:"base_url,optional"`
	Description string   `hcl:"description,optional"`
}

type arg struct {
	Name string `cty:"name"`
	Type string `cty:"type"`
}

func (h *hclDatasource) List() map[string]errorDefinition {
	if !h.isLoaded() {
		h.load()
	}

	return h.list
}

func (h *hclDatasource) isLoaded() bool {
	h.Lock()
	defer h.Unlock()

	return h.list != nil
}

func (h *hclDatasource) load() {
	if h.isLoaded() {
		return
	}

	h.Lock()
	defer h.Unlock()

	h.list = make(map[string]errorDefinition, len(h.Errors))

	for _, e := range h.Errors {
		h.list[e.Code] = errorDefinition{
			Code:       e.Code,
			Message:    e.Message,
			Cause:      e.Cause,
			Guide:      e.Guide,
			Args:       h.argsMap(e),
			Categories: e.Categories,
			Labels:     e.Labels,
		}
	}
}

func (h *hclDatasource) argsMap(e hclErrorDefinition) map[string]arg {
	if len(e.Args) <= 0 {
		return nil
	}

	args := make(map[string]arg, len(e.Args))

	for _, argRaw := range e.Args {
		var arg arg
		_ = gocty.FromCtyValue(argRaw, &arg)
		args[arg.Name] = arg
	}
	return args
}

func (h *hclDatasource) Options() errorOptions {
	if h.Opts == nil {
		return errorOptions{}
	}

	return *h.Opts
}

func (h *hclDatasource) FindByCode(code string) (errorDefinition, bool) {
	err, ok := h.list[code]
	if !ok {
		// if we cannot find the error by code, create one
		return errorDefinition{
			Code: code,
		}, false
	}
	return err, true
}

func (h *hclDatasource) Validate() error {
	for _, e := range h.list {
		for k := range e.Labels {
			if _, found := e.Args[k]; found {
				return NewInvalidDefinitionsErr(NewArgumentLabelNameClashErr(nil, k), h.path)
			}
		}
	}

	return nil
}

func (h *hclDatasource) Hash() string {
	return fmt.Sprintf("%x", md5.Sum(h.source))
}

func (h *hclDatasource) SchemaVersion() string {
	return h.Version
}

func NewHCLDatasource(path string) (DataSource, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, NewFileNotFoundErr(err, path)
	}

	db, err := parseHCL(path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func parseHCL(path string) (*hclDatasource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, NewFileNotReadableErr(err, path)
	}

	b, err := io.ReadAll(f)
	if err != nil || len(b) == 0 {
		return nil, NewInvalidDatasourceErr(err, path)
	}

	var db hclDatasource

	err = hclsimple.Decode(filepath.Base(path), b, &hcl.EvalContext{
		Functions: map[string]function.Function{
			"arg": function.New(&function.Spec{
				Params: []function.Parameter{
					{Name: "name", Type: cty.String, AllowNull: false, AllowDynamicType: false},
					{Name: "type", Type: cty.String, AllowNull: false, AllowDynamicType: false},
				},
				Type: function.StaticReturnType(cty.Object(map[string]cty.Type{
					"name": cty.String,
					"type": cty.String,
				})),
				Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
					return cty.ObjectVal(map[string]cty.Value{
						"name": args[0],
						"type": args[1],
					}), nil
				},
			}),
			"file": function.New(&function.Spec{
				Params: []function.Parameter{
					{Name: "path", Type: cty.String, AllowNull: false, AllowDynamicType: false},
				},
				Type: function.StaticReturnType(cty.String),
				Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
					path := filepath.Clean(args[0].AsString())
					guide, err := guideFromFile(path)
					return cty.StringVal(guide), err
				},
			}),
		},
	}, &db)

	if err != nil {
		return nil, NewInvalidSyntaxErr(err, path)
	}

	db.load()
	db.source = b
	db.path = path
	return &db, nil
}

func guideFromFile(path string) (string, error) {
	_, err := os.Open(path)
	if err != nil {
		return "", NewFileNotFoundErr(err, path)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", NewFileNotReadableErr(err, path)
	}

	return string(b), nil
}
