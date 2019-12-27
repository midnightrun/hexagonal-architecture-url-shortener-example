package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case shortener.Info:
		return zapcore.InfoLevel
	case shortener.Warn:
		return zapcore.WarnLevel
	case shortener.Debug:
		return zapcore.DebugLevel
	case shortener.Error:
		return zapcore.ErrorLevel
	case shortener.Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config shortener.Configuration) (logger shortener.Logger, error) {

}
