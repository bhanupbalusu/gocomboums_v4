package logs

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func init() {
	Logger = log.New()
	Logger.SetFormatter(&log.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(log.InfoLevel) // Change this level based on your need
}

// Debug logs a message at level Debug
func Debug(msg string) {
	Logger.Debug(msg)
}

// Debugf logs a formatted message at level Debug
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info logs a message at level Info
func Info(msg string) {
	Logger.Info(msg)
}

// Infof logs a formatted message at level Info
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn logs a message at level Warn
func Warn(msg string) {
	Logger.Warn(msg)
}

// Warnf logs a formatted message at level Warn
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error logs a message at level Error
func Error(msg string, err ...error) {
	if len(err) > 0 {
		Logger.Error(fmt.Sprintf("%s: %v", msg, err[0]))
	} else {
		Logger.Error(msg)
	}
}

// Errorf logs a formatted message at level Error
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(msg string) {
	Logger.Fatal(msg)
}

// Fatalf logs a formatted message at level Fatal then the process will exit with status set to 1
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panic logs a message at level Panic, then panics
func Panic(msg string) {
	Logger.Panic(msg)
}

// Panicf logs a formatted message at level Panic, then panics
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// WithFields provides the ability to easily add multiple fields to a log entry
func WithFields(fields log.Fields) *log.Entry {
	return Logger.WithFields(fields)
}
