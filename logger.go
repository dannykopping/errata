package errata

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type LogLevel interface {
	Erratum
	GetLogLevel() string
}

func LogError(logger log.Logger, err error) {
	lvl := level.ErrorValue()

	if e, ok := err.(LogLevel); ok {
		// see if the error we received implements the LogLevel interface, and use that level
		lvl = level.ParseDefault(e.GetLogLevel(), lvl)
	}

	if e, ok := err.(Erratum); ok {
		// see if the error we received implements the Erratum interface, and enhance the output with that data
		logger = log.With(logger, "code", e.Code(), "msg", e.Message(), "uuid", e.UUID(), "help", e.HelpURL())
	} else {
		// if this is not an Erratum, just use the error message
		logger = log.With(logger, "err", err.Error())
	}

	log.WithPrefix(logger, level.Key(), lvl).Log()
}
