package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
}

func New() *Logger {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

	return &Logger{
		debug: log.New(os.Stdout, "DEBUG: ", flags),
		info:  log.New(os.Stdout, "INFO: ", flags),
		error: log.New(os.Stderr, "ERROR: ", flags),
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.debug.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.error.Output(2, fmt.Sprintf(format, v...))
}