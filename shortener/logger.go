package shortener

import "errors"

var log Logger

type Fields map[string]interface{}

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstace = errors.New("invalid logger instance")
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	WithFields(keyValues Fields)
}

type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLovation      string
}
