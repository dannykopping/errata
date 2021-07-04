package errata

type NullDataSource struct {}

func NewNullDataSource() *NullDataSource {
	return &NullDataSource{}
}

func (n *NullDataSource) FindByCode(code string) Error {
	return Error{
		Code: code,
	}
}
