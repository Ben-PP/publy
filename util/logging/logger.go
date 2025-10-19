package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/DeRuina/timberjack"
)

// GetLogger initializes and returns a configured slog.Logger instance.
func GetLogger() *slog.Logger {
	logDir := "/var/log/publy"
	if os.Getenv("GO_ENV") == "dev" {
		logDir = "."
	}
	_, err := os.Stat(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			panic("Log directory does not exist.")
		} else {
			panic("Unknown error while testing log directory.")
		}
	}
	testFile := fmt.Sprintf("%s/test.log", logDir)
	f, err := os.Create(testFile)
	if err != nil {
		panic(fmt.Sprintf("Could not write the log directory '%s'", logDir))
	}
	f.Close()
	os.Remove(testFile)

	logRotator := &timberjack.Logger{
		Filename:   fmt.Sprintf("%s/app.log", logDir),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     100,
		Compress:   true,
	}

	appMultiWriter := slog.NewJSONHandler(
		io.MultiWriter(logRotator),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		},
	)

	appLogger := slog.New(appMultiWriter).With(slog.String("program_name", "Publy"))

	return appLogger
}

// log is a helper function to log messages with a specific sourcetype.
func log(level slog.Leveler, message string, sourcetype string, args ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		level.Level(),
		message,
		append(
			[]slog.Attr{slog.String("sourcetype", "publy:"+sourcetype)},
			args...,
		)...,
	)
}
