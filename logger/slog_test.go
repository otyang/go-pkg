package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNewSlogLogger(t *testing.T) {
	addSource := true
	output := os.Stdout
	level := LevelDebug

	// expected
	handler := slog.HandlerOptions{AddSource: addSource, Level: slog.LevelDebug}
	jsonHandler := slog.NewJSONHandler(output, &handler)
	expected1 := &SlogLogger{logger: slog.New(jsonHandler), reset: slog.New(jsonHandler)}

	// actual
	actual1 := NewSlogLogger(level, DisplayTypeJson, addSource, output)

	// test 1
	assert.Equal(t, expected1, actual1, "test case 1: it should be same. but isnt, why?")

	// test 2
	actual1.With("testing-service", "hosted-on-google-cloud")

	actual1.Reset()
	assert.Equal(t, expected1, actual1, "test case 2: it should be same. but isnt, why?")
}
