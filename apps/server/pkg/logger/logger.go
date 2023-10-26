package logger

import (
	"github.com/chat-system/server/pkg/config"
	"go.uber.org/zap"
)

var defaultLogger Logger

func Init(conf *config.LoggerConfig) {
	// parse zap level
	zapLevel, err := zap.ParseAtomicLevel(conf.Level)

	if err != nil {
		return
	}

	// initialize from config
	zapConfig := zap.Config{
		Level:            zapLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zap, err := zapConfig.Build()

	if err != nil {
		return
	}

	defaultLogger = Logger{
		zap: zap.Sugar(),
	}
}

func Infow(msg string, keyAndValues ...interface{}) {
	if defaultLogger.zap == nil {
		return
	}

	defaultLogger.Infow(msg, keyAndValues...)
}

func Debugw(msg string, keyAndValues ...interface{}) {
	if defaultLogger.zap == nil {
		return
	}

	defaultLogger.Debugw(msg, keyAndValues...)
}

func Warnw(msg string, err error, keyAndValues ...interface{}) {
	if defaultLogger.zap == nil {
		return
	}

	defaultLogger.Warnw(msg, err, keyAndValues...)
}

func Errorw(msg string, err error, keyAndValues ...interface{}) {
	if defaultLogger.zap == nil {
		return
	}

	defaultLogger.Errorw(msg, err, keyAndValues...)
}

type Logger struct {
	zap *zap.SugaredLogger
}

func (l *Logger) Debugw(msg string, keyAndValues ...interface{}) {
	l.zap.Infow(msg, keyAndValues...)
}

func (l *Logger) Infow(msg string, keyAndValues ...interface{}) {
	l.zap.Infow(msg, keyAndValues...)
}

func (l *Logger) Warnw(msg string, err error, keyAndValues ...interface{}) {
	if err != nil {
		keyAndValues = append(keyAndValues, "error", err)
	}

	l.zap.Warnw(msg, keyAndValues...)
}

func (l *Logger) Errorw(msg string, err error, keyAndValues ...interface{}) {
	if err != nil {
		keyAndValues = append(keyAndValues, "error", err)
	}

	l.zap.Errorw(msg, keyAndValues...)
}
