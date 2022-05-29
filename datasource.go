package errata

type DataSource interface {
	List() map[string]ErrorDefinition
	FindByCode(code string) ErrorDefinition
	// Version must indicate a unique version based on the given source data
	Version() string
}
