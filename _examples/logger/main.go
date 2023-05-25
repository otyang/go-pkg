package main

import (
	"fmt"
	"os"

	"github.com/otyang/go-pkg/logger"
)

func main() {
	var (
		logLevel             = logger.LevelDebug      // logger.LevelWarn, logger.LevelInfo, logger.LevelError
		displayJsonOrText    = logger.DisplayTypeJson // logger.DisplayTypeText
		addSourceCodeAndLine = false
		writeTo              = os.Stdout
	)

	// logger via New Golang SLOG: golang.org/x/exp/slog
	var log logger.Interface = logger.NewSlogLogger(logLevel, displayJsonOrText, addSourceCodeAndLine, writeTo)

	fmt.Println("	-----------0		")
	// Log at info level
	log.Info("Hello world at info level")

	// Log With Attributes. This attaches this attributes to every log
	log.With(
		"appName", "helloWorld micro-service",
		"appAddress", "https://helloworld.com:80",
	)

	log.Info("Hello world at info level") // Now notice the difference between this and first example.
	fmt.Println("	-----------1		")

	// Helper func that calls log.With() and adds the binary build info
	logger.WithBuildInfo(log)

	log.Reset() // resets the logger to way it was instantiated removing the attributes

	log.Info("cant connect to cache")
	log.Info("cant connect to cache", "requestId", "1")
	fmt.Println("	-----------2		")
	log.Debug("cant connect to cache")
	log.Debug("cant connect to cache", "requestId", "1cdhg")
	fmt.Println("	-----------3		")
	log.Fatal("cant connect to cache")
	log.Fatal("cant connect to cache", "requestId", "1cdhg")
	fmt.Println("	-----------4		")
	log.Warn("cant connect to cache")
	log.Warn("cant connect to cache", "requestId", "1cdhg")
	fmt.Println("	-----------		")
	log.Error("cant connect to cache")
	log.Error("cant connect to cache", "requestId", "1cdhg")
}
