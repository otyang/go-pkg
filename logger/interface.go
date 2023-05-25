package logger

import (
	"os"
	"runtime/debug"
)

// Interface is a logger interface that ensures different log library can be used and swap with ease.
type Interface interface {
	// With attached attributes to eveery log
	With(args ...any)
	// Resets the logger to how it was instantiated.
	Reset()
	// Debug log a message at DEBUG level
	Debug(message string, args ...any)
	// Error log a message at ERROR level
	Error(message string, args ...any)
	// Fatal log a message at FATAL level
	Fatal(message string, args ...any)
	// Info log a message at INFO level
	Info(message string, args ...any)
	// Warn log a message at WARN level
	Warn(message string, args ...any)
}

// LogLevels define the various levels of log available
const (
	LevelError LogLevel = "error"
	LevelWarn  LogLevel = "warn"
	LevelInfo  LogLevel = "info"
	LevelDebug LogLevel = "debug"
)

type LogLevel string

func (l LogLevel) String() string {
	return string(l)
}

func (l LogLevel) IsValid() bool {
	switch l {
	case LevelError, LevelWarn, LevelInfo, LevelDebug:
		return true
	default:
		return false
	}
}

const (
	DisplayTypeJson LogDisplayType = "json"
	DisplayTypeText LogDisplayType = "text"
)

type LogDisplayType string

func (l LogDisplayType) String() string {
	return string(l)
}

func (l LogDisplayType) IsValid() bool {
	switch l {
	case DisplayTypeJson, DisplayTypeText:
		return true
	default:
		return false
	}
}

// WithBuildInfo is an utility func that attaches the cmd build information -
// *	process id of the caller,
// *	golang version,
// *	args other params
func WithBuildInfo(logger Interface) {
	osGetPID := os.Getpid()
	buildInfo, _ := debug.ReadBuildInfo()

	logger.With(
		"program-pid", osGetPID,
		"go-version", buildInfo.GoVersion,
	)
}
