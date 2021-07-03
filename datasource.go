package errata

type NoopDataSource struct {}

func NewNoopDataSource() *NoopDataSource {
	return &NoopDataSource{}
}

func (n *NoopDataSource) FindByCode(code string) error {
	return New(code)
}

