package errata

type DataSource interface {
	List() map[string]ErrorDefinition
	Options() ErrorOptions
	FindByCode(code string) ErrorDefinition
	Validate() error
	// Version must indicate a unique version based on the given source data
	Version() string
}
