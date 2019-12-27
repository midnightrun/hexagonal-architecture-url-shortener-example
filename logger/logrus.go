package logrus

import (
	"io"
	"os"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{}
	}

	return &logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
}

func NewLogrusLogger(config shortener.Configuration) (shortener.Logger, error) {
	logLevel := config.ConsoleLevel

	if logLevel == "" {
		logLevel = config.FileLevel
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	stdOutHandler := os.Stdout

	fileHandler := &lumberjack.Logger{
		Filename: config.FileLocation,
		MaxSize:  100,
		Compress: true,
		MaxAge:   28,
	}

	lLogger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(config.ConsoleJSONFormat),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if config.EnableConsole && config.EnableFile {
		lLogger.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))
	} else {
		lLogger.SetOutput(fileHandler)
		lLogger.SetFormatter(getFormatter(config.FileJSONFormat))
	}

	return &logrusLogger{
		logger: lLogger,
	}, nil
}

func convertToLogrusFields(fields shortener.Fields) logrus.Fields {
	logrusFields := logrus.Fields{}

	for k, v := range fields {
		logrusFields[k] = v
	}

	return logrusFields
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) WithFields(keyValues shortener.Fields) shortener.Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(keyValues)),
	}
}

func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) WithFields(keyValues shortener.Fields) shortener.Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(keyValues)),
	}
}
