package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func New() *Logger {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

	return &Logger{
		debug: log.New(os.Stdout, "DEBUG: ", flags),
		info:  log.New(os.Stdout, "INFO: ", flags),
		warn:  log.New(os.Stdout, "WARN: ", flags),
		error: log.New(os.Stderr, "ERROR: ", flags),
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	_ = l.debug.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(format string, v ...interface{}) {
	_ = l.info.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(format string, v ...interface{}) {
	_ = l.warn.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	_ = l.error.Output(2, fmt.Sprintf(format, v...))
}

var defaultLogger = New()

// GetLogger возвращает глобальный экземпляр логгера
func GetLogger() *Logger {
	return defaultLogger
}
