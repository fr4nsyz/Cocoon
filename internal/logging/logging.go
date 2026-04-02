package logging

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger  *log.Logger
	verbose bool
}

func New(verbose bool) *Logger {
	return &Logger{
		logger:  log.New(os.Stdout, "", 0),
		verbose: verbose,
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.logger.Printf("[WARN] "+format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.logger.Printf("[ERROR] "+format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	if l.verbose {
		l.logger.Printf("[DEBUG] "+format, args...)
	}
}

func (l *Logger) Blocked(action string, details string) {
	msg := fmt.Sprintf("BLOCKED %s: %s", action, details)
	if l.verbose {
		l.logger.Printf("[BLOCKED] %s", msg)
	}
}
