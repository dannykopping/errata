package errata

type ErrorDefinition struct {
	Code       string            `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	Cause      string            `json:"cause,omitempty"`
	Guide      string            `json:"guide,omitempty"`
	Args       []Arg             `json:"args,omitempty"`
	Categories []string          `json:"categories,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type ErrorOptions struct {
	Prefix      string   `hcl:"prefix,optional"`
	Imports     []string `hcl:"imports,optional"`
	BaseURL     string   `hcl:"base_url,optional"`
	Description string   `hcl:"description,optional"`
}

type Arg struct {
	Name string `cty:"name"`
	Type string `cty:"type"`
}
