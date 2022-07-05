package errata

import (
	"errors"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type LogLevel interface {
	Erratum
	GetLogLevel() string
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
}

func LogError(err error, args ...interface{}) {
	lvl := level.ErrorValue()

	if e, ok := err.(LogLevel); ok {
		// see if the error we received implements the LogLevel interface, and use that level
		lvl = level.ParseDefault(e.GetLogLevel(), lvl)
	}

	// create a temporary logger so we don't overwrite the global logger
	var l log.Logger
	l = log.With(logger, args...)

	var parentUUID string

	if e, ok := err.(Erratum); ok {
		parentUUID = e.UUID()

		// see if the error we received implements the Erratum interface, and enhance the output with that data
		l = log.With(l, "code", e.Code(), "err", e.Error(), "uuid", e.UUID(), "help", e.HelpURL())
		// add arguments if any have been defined
		l = log.With(l, argsToKeyVals(e.Args())...)
	} else {
		// if this is not an Erratum, just use the error message
		l = log.With(l, "err", err.Error())
	}

	log.WithPrefix(l, level.Key(), lvl).Log()

	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		LogError(unwrapped, "parent-uuid", parentUUID)
	}
}

func argsToKeyVals(args map[string]interface{}) []interface{} {
	var kvs []interface{}
	for k, v := range args {
		kvs = append(kvs, k, v)
	}
	return kvs
}
