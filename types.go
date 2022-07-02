package errata

type errorDefinition struct {
	Code       string
	Message    string
	Cause      string
	Guide      string
	Categories []string
	Args       []arg
	Labels     map[string]string
}
