package errata

type ErrorDefinition struct {
	Code       string
	Message    string
	Cause      string
	Guide      string
	Args       []Arg
	Categories []string
	Labels     map[string]string
}

type ErrorOptions struct {
	Prefix     string   `hcl:"prefix"`
	GuidesPath string   `hcl:"guides_path"`
	Imports    []string `hcl:"imports"`
}

type Arg struct {
	Name string `cty:"name"`
	Type string `cty:"type"`
}
