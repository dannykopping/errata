package errata

type DataSource interface {
	List() map[string]errorDefinition
	Options() errorOptions
	// FindByCode looks up an Erratum by a given code and is guaranteed to return an errorDefinition.
	// However, if the code cannot be found, an errorDefinition will be crated with just the given code
	// and false will be returned as the second return value
	FindByCode(code string) (errorDefinition, bool)
	Validate() error
	// Hash returns the hashed contents of the given datasource, to detect drift
	Hash() string
	// SchemaVersion returns the schema version of the given datasource
	SchemaVersion() string
}
