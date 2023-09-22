package logger

import (
	"github.com/Hvaekar/med-account/config"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func newLogrusLogger() Logger {
	return &logrusLogger{}
}

var logrusLoggerLevel = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"trace": logrus.TraceLevel,
}

func (l *logrusLogger) getLevel(level string) logrus.Level {
	lvl, ok := logrusLoggerLevel[level]
	if !ok {
		return logrus.DebugLevel
	}

	return lvl
}

func (l *logrusLogger) Init(cfg *config.Config) error {
	level := l.getLevel(cfg.Logger.Level)

	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: level,
	}

	if cfg.Logger.Encoding == "json" {
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	} else {
		logger.Formatter = &logrus.TextFormatter{
			ForceColors:     true,
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		}
	}

	if !cfg.Logger.DisableCaller {
		logger.ReportCaller = true
	}

	l.logger = &logger

	return nil
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogger{entry: l.logger.WithFields(convertToLogrusFields(fields))}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for k, v := range fields {
		logrusFields[k] = v
	}

	return logrusFields
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.entry != nil {
		l.entry.Debug(args...)
	} else {
		l.logger.Debug(args...)
	}
}

func (l *logrusLogger) Debugf(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Debugf(template, args...)
	} else {
		l.logger.Debugf(template, args...)
	}
}

func (l *logrusLogger) Info(args ...interface{}) {
	if l.entry != nil {
		l.entry.Info(args...)
	} else {
		l.logger.Info(args...)
	}
}

func (l *logrusLogger) Infof(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Infof(template, args...)
	} else {
		l.logger.Infof(template, args...)
	}
}

func (l *logrusLogger) Warn(args ...interface{}) {
	if l.entry != nil {
		l.entry.Warn(args...)
	} else {
		l.logger.Warn(args...)
	}
}

func (l *logrusLogger) Warnf(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Warnf(template, args...)
	} else {
		l.logger.Warnf(template, args...)
	}
}

func (l *logrusLogger) Error(args ...interface{}) {
	if l.entry != nil {
		l.entry.Error(args...)
	} else {
		l.logger.Error(args...)
	}
}

func (l *logrusLogger) Errorf(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Errorf(template, args...)
	} else {
		l.logger.Errorf(template, args...)
	}
}

func (l *logrusLogger) Panic(args ...interface{}) {
	if l.entry != nil {
		l.entry.Panic(args...)
	} else {
		l.logger.Panic(args...)
	}
}

func (l *logrusLogger) Panicf(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Panicf(template, args...)
	} else {
		l.logger.Panicf(template, args...)
	}
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatal(args...)
	} else {
		l.logger.Fatal(args...)
	}
}

func (l *logrusLogger) Fatalf(template string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Debugf(template, args...)
	} else {
		l.logger.Debugf(template, args...)
	}
}

func (l *logrusLogger) IsLevel(lvl string) bool {
	if l.entry != nil {
		return logrusLoggerLevel[lvl] >= l.entry.Level
	} else {
		return logrusLoggerLevel[lvl] >= l.logger.Level
	}
}
