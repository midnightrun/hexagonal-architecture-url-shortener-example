package shortener

type Fields map[string]interface{}

const (
	Debug = "debug"
	Info  = "info"
	Warn  = "warn"
	Error = "error"
	Fatal = "fatal"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
	ConsoleLevel      string
}
