package go_bagit

import (
	"io"
	stdlog "log"
)

var log = stdlog.Default()

func init() {
	log.SetOutput(io.Discard)
}

// WithLogger replaces the logger.
func WithLogger(logger *stdlog.Logger) {
	if logger == nil {
		panic("WithLogger expects a non-nil value")
	}
	log = logger
}

// Logger returns the current logger.
func Logger() *stdlog.Logger {
	return log
}
