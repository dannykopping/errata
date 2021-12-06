package errors

import "fmt"

type Error struct {
	Code       string            `yaml:"code"`
	Message    string            `yaml:"message"`
	Cause      string            `yaml:"cause,omitempty"`
	Categories []string          `yaml:"categories,omitempty,flow"`
	Labels     map[string]string `yaml:"labels,omitempty,flow"`

	HTTP  HTTP  `yaml:"http,omitempty"`
	Shell Shell `yaml:"shell,omitempty"`

	inner error
}

type HTTP struct {
	StatusCode int `yaml:"status_code"`
}

type Shell struct {
	ExitCode int `yaml:"exit_code"`
}

func (e Error) Wrap(err error) error {
	if err == nil {
		return nil
	}

	e.inner = err
	return e
}

func (e Error) Unwrap() error {
	return e.inner
}

func (e Error) Error() string {
	if w := e.Unwrap(); w != nil {
		return fmt.Sprintf("[%s] %s", e.Code, w)
	}

	return e.Code
}
