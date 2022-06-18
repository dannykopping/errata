package errata

type errorDefinition struct {
	Code       string
	Message    string
	Cause      string
	Guide      string
	Args       []arg
	Categories []string
	Labels     map[string]string
}
