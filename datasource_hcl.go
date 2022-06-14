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
	sync.RWMutex

	source []byte
	list   map[string]ErrorDefinition

	SchemaVersion string               `hcl:"version"`
	Opts          ErrorOptions         `hcl:"options,block"`
	Errors        []hclErrorDefinition `hcl:"error,block"`
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

func (h *hclDatasource) List() map[string]ErrorDefinition {
	if !h.isLoaded() {
		h.load()
	}

	return h.list
}

func (h *hclDatasource) isLoaded() bool {
	h.RLock()
	defer h.RUnlock()

	return h.list != nil
}

func (h *hclDatasource) load() {
	if h.isLoaded() {
		return
	}

	h.Lock()
	defer h.Unlock()

	h.list = make(map[string]ErrorDefinition, len(h.Errors))

	for _, e := range h.Errors {
		var args []Arg
		for _, argRaw := range e.Args {
			var arg Arg
			_ = gocty.FromCtyValue(argRaw, &arg)
			args = append(args, arg)
		}

		h.list[e.Code] = ErrorDefinition{
			Code:       e.Code,
			Message:    e.Message,
			Cause:      e.Cause,
			Guide:      e.Guide,
			Args:       args,
			Categories: e.Categories,
			Labels:     e.Labels,
		}
	}
}

func (h *hclDatasource) Options() ErrorOptions {
	return h.Opts
}

func (h *hclDatasource) FindByCode(code string) (ErrorDefinition, bool) {
	err, ok := h.list[code]
	if !ok {
		// if we cannot find the error by code, create one
		return ErrorDefinition{
			Code: code,
		}, false
	}
	return err, true
}

func (h *hclDatasource) Validate() error {
	// TODO: define validation rules
	return nil
}

func (h *hclDatasource) Version() string {
	return fmt.Sprintf("%x", md5.Sum(h.source))
}

func NewHCLDatasource(path string) (DataSource, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, NewFileNotFoundErr(err, path)
	}

	db, err := parseHCL(path)
	if err != nil {
		return nil, NewInvalidDatasourceErr(err)
	}

	return db, nil
}

func parseHCL(path string) (*hclDatasource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, NewFileNotReadableErr(err, path)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
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
					return cty.StringVal(fmt.Sprintf("file://%s", path)), err
				},
			}),
		},
	}, &db)

	if err != nil {
		return nil, NewInvalidSyntaxErr(err)
	}

	db.load()
	db.source = b
	return &db, nil
}
