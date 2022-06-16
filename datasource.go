package errata

type DataSource interface {
	List() map[string]ErrorDefinition
	Options() ErrorOptions
	// FindByCode looks up an erratum by a given code and is guaranteed to return an ErrorDefinition.
	// However, if the code cannot be found, an ErrorDefinition will be crated with just the given code
	// and false will be returned as the second return value
	FindByCode(code string) (ErrorDefinition, bool)
	Validate() error
	// Version must indicate a unique version based on the given source data
	Version() string
}
