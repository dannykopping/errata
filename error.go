package errata

type Error struct {
	Code       string	`yaml:"code"`
	Message    string   `yaml:"message"`
	Cause      string   `yaml:"cause,omitempty"`
	Categories []string `yaml:"categories,omitempty,flow"`

	External *Error `yaml:"external,omitempty"`
	HTTP     *HTTP  `yaml:"http,omitempty"`
}

func (e *Error) Error() string {
	return e.Code
}