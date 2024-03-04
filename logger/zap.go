package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap Logger -.
type Zap struct {
	logger *zap.Logger
}

var _ Logger = (*Zap)(nil)

// New -.
func NewZap(level string) *Zap {
	var cfg zap.Config

	switch strings.ToLower(level) {
	case "error":
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zapcore.ErrorLevel)
	case "warn":
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zapcore.WarnLevel)
	case "info":
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zapcore.InfoLevel)
	case "debug":
		cfg = zap.NewDevelopmentConfig()
	default:
		cfg = zap.NewProductionConfig()
		cfg.Level.SetLevel(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = logger.WithOptions(zap.AddCallerSkip(3))

	return &Zap{
		logger: logger,
	}
}

// Debug -.
func (l *Zap) Debug(message interface{}, args ...interface{}) {
	l.sugar().Debugw(fmt.Sprint(message), args...)
}

// Info -.
func (l *Zap) Info(message string, args ...interface{}) {
	l.sugar().Infow(message, args...)
}

// Warn -.
func (l *Zap) Warn(message string, args ...interface{}) {
	l.sugar().Warnw(message, args...)
}

// Error -.
func (l *Zap) Error(message interface{}, args ...interface{}) {
	l.sugar().Errorw(fmt.Sprint(message), args...)
}

// Fatal -.
func (l *Zap) Fatal(message interface{}, args ...interface{}) {
	l.sugar().Fatalw(fmt.Sprint(message), args...)
}

// sugar returns a zap.SugaredLogger for convenience
func (l *Zap) sugar() *zap.SugaredLogger {
	return l.logger.Sugar()
}
