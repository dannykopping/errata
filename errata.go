package errata

var (
	source DataSource = NewNullDataSource()
)

var (
	InvalidDataSource = Error{
		Code:    "invalid_datasource",
		Message: "Given Errata datasource is invalid",
	}
)

func RegisterDataSource(ds DataSource) error {
	if ds == nil {
		return InvalidDataSource
	}

	source = ds
	return nil
}

func New(code string) Error {
	return source.FindByCode(code)
}
