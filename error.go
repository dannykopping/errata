package errata

type Error struct {
	Code       string	`yaml:"code"`
	Message    string   `yaml:"message"`
	Cause      string   `yaml:"cause,omitempty"`
	Categories []string `yaml:"categories,omitempty,flow"`

	External *Error `yaml:"external,omitempty"`
	HTTP     *HTTP  `yaml:"http,omitempty"`
	Shell    *Shell `yaml:"shell,omitempty"`
}

func (e Error) Error() string {
	return e.Code
}

func (e Error) HTTPStatusCode(defaultCode int) int {
	if e.HTTP == nil || e.HTTP.StatusCode == 0 {
		return defaultCode
	}

	return e.HTTP.StatusCode
}

func (e Error) ShellExitCode(defaultCode int) int {
	if e.Shell == nil || e.Shell.ExitCode == 0 {
		return defaultCode
	}

	return e.Shell.ExitCode
}