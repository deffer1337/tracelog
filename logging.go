package tracelog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	NormalLogLevel = "NORMAL"
	DevelLogLevel  = "DEVEL"
	ErrorLogLevel  = "ERROR"
	timeFlags      = log.LstdFlags | log.Lmicroseconds
)

var InfoLogger = NewErrorLogger(os.Stdout, "INFO: ")
var WarningLogger = NewErrorLogger(os.Stdout, "WARNING: ")
var ErrorLogger = NewErrorLogger(os.Stderr, "ERROR: ")
var DebugLogger = NewErrorLogger(ioutil.Discard, "DEBUG: ")

var LogLevels = []string{NormalLogLevel, DevelLogLevel, ErrorLogLevel}
var logLevel = NormalLogLevel
var logLevelFormatters = map[string]string{
	NormalLogLevel: "%v",
	ErrorLogLevel:  "%v",
	DevelLogLevel:  "%+v",
}

func setupLoggers() {
	if logLevel == NormalLogLevel {
		DebugLogger = NewErrorLogger(ioutil.Discard, "DEBUG: ")
	} else if logLevel == ErrorLogLevel {
		DebugLogger = NewErrorLogger(ioutil.Discard, "DEBUG: ")
		InfoLogger = NewErrorLogger(ioutil.Discard, "INFO: ")
		WarningLogger = NewErrorLogger(ioutil.Discard, "WARNING: ")
	} else {
		DebugLogger = NewErrorLogger(os.Stdout, "DEBUG: ")
	}
}

type LogLevelError struct {
	error
}

func NewLogLevelError(incorrectLogLevel string) LogLevelError {
	return LogLevelError{errors.Errorf("got incorrect log level: '%s', expected one of: '%v'", incorrectLogLevel, LogLevels)}
}

func (err LogLevelError) Error() string {
	return fmt.Sprintf(GetErrorFormatter(), err.error)
}

func GetErrorFormatter() string {
	return logLevelFormatters[logLevel]
}

func UpdateLogLevel(newLevel string) error {
	isCorrect := false
	for _, level := range LogLevels {
		if newLevel == level {
			isCorrect = true
		}
	}
	if !isCorrect {
		return NewLogLevelError(newLevel)
	}

	logLevel = newLevel
	setupLoggers()
	return nil
}
