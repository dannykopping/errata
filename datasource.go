package errata

type DataSource interface {
	FindByCode(code string) Error
}

var (
	DatasourceUninitializedError = Error{
		Code:    "datasource_uninitialized",
		Message: "Configured database could not be initialized",
	}
)

