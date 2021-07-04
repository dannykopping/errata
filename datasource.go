package errata

type NoopDataSource struct {}

func NewNoopDataSource() *NoopDataSource {
	return &NoopDataSource{}
}

func (n *NoopDataSource) FindByCode(code string) Error {
	return New(code)
}

