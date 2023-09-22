package logger

import (
	"github.com/Hvaekar/med-account/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func newZapLogger() Logger {
	return &zapLogger{}
}

var zapLoggerLevel = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

func (l *zapLogger) getLevel(level string) zapcore.Level {
	lvl, ok := zapLoggerLevel[level]
	if !ok {
		return zapcore.DebugLevel
	}

	return lvl
}

func (l *zapLogger) Init(cfg *config.Config) error {
	level := l.getLevel(cfg.Logger.Level)

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Server.Development,
		DisableCaller:     cfg.Logger.DisableCaller,
		DisableStacktrace: cfg.Logger.DisableStacktrace,
		Encoding:          cfg.Logger.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "ts",
			EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
		},
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: make(map[string]interface{}),
	}

	log, err := zapConfig.Build()
	if err != nil {
		return err
	}

	l.logger = log.Sugar()

	return nil
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	f := make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}

	log := l.logger.With(f...)

	return &zapLogger{logger: log}
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *zapLogger) IsLevel(lvl string) bool {
	return zapLoggerLevel[lvl] >= l.logger.Level()
}
