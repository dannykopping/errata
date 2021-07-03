package errata

type DataSource interface {
	FindByCode(code string) error
}

type Metadata struct {
	Key   string
	Value interface{}
}

var (
	source DataSource = NewNoopDataSource()
)

var (
	InvalidDataSource = Error{
		Code:    "invalid_datasource",
		Message: "Given Errata datasource is invalid",
	}
)

func RegisterSource(ds DataSource) error {
	if ds == nil {
		return &InvalidDataSource
	}

	source = ds
	return nil
}

func New(code string) error {
	return source.FindByCode(code)
}
