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
	Prefix  string   `hcl:"prefix,optional"`
	Imports []string `hcl:"imports,optional"`
	BaseURL string   `hcl:"base_url,optional"`
}

type Arg struct {
	Name string `cty:"name"`
	Type string `cty:"type"`
}
