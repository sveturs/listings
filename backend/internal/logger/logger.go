// Package logger
// backend/internal/logger/logger.go
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Init(env string, level string) error {
	if level == "" {
		level = "info"
	}
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	if env == "dev" || env == "development" {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(l)
	return nil
}

func Get() *zerolog.Logger {
	return &log
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Debug() *zerolog.Event {
	return log.Debug()
}
