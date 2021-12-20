package errata

type DataSource interface {
	List() map[string]ErrorDefinition
	FindByCode(code string) ErrorDefinition
}
