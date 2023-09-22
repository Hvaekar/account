package logger

import "github.com/Hvaekar/med-account/config"

type Fields map[string]interface{}

type Logger interface {
	Init(cfg *config.Config) error
	WithFields(fields Fields) Logger
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	IsLevel(lvl string) bool
}

func Get(loggerName string) Logger {
	switch loggerName {
	case "zap":
		return newZapLogger()
	case "logrus":
		return newLogrusLogger()
	}

	return nil
}
