package logger

import (
	"log"
	"os"
)

type Logger struct {
	Verbose bool
	logger  *log.Logger
}

func NewLogger(verbose bool) *Logger {
	return &Logger{
		Verbose: verbose,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	if l.Verbose {
		l.logger.Println("[INFO] ", msg)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Verbose {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Error(msg string) {
	l.logger.Println("[ERROR] ", msg)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Fatal(msg string) {
	l.logger.Fatalf("[FATAL] %s", msg)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, v...)
}
