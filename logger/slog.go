package logger

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
)

var _ Interface = (*SlogLogger)(nil)

// SlogLogger wraps the sirupsen/logrus library in a way that is compatible with our cache interface
type SlogLogger struct {
	logger *slog.Logger

	reset *slog.Logger // for reset purpose
}

func NewSlogLogger(level LogLevel, logFormatType LogDisplayType, AddSource bool, output io.Writer) *SlogLogger {
	var logLevel slog.Level

	switch level {
	case LevelDebug:
		logLevel = slog.LevelDebug
	case LevelInfo:
		logLevel = slog.LevelInfo
	case LevelWarn:
		logLevel = slog.LevelWarn
	case LevelError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	handler := slog.HandlerOptions{AddSource: AddSource, Level: logLevel}

	jsonHandler := slog.NewJSONHandler(output, &handler)
	textHandler := slog.NewTextHandler(output, &handler)

	jsonLog := slog.New(jsonHandler)
	textLog := slog.New(textHandler)
	resetJsonLog := slog.New(jsonHandler)
	resetTextLog := slog.New(textHandler)

	switch logFormatType {
	case DisplayTypeJson:
		return &SlogLogger{logger: jsonLog, reset: resetJsonLog}
	default:
		return &SlogLogger{logger: textLog, reset: resetTextLog}
	}
}

func (l *SlogLogger) With(args ...any) {
	l.logger = l.logger.With(args...)
}

func (l *SlogLogger) Reset() {
	l.logger = l.reset
}

func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) Fatal(msg string, args ...any) {
	l.logger.Warn(msg, args...)
	os.Exit(1)
}

func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
