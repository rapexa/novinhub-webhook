package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// New creates a new logger instance
func New() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &Logger{Logger: logger}
}

// NewWithLevel creates a new logger instance with specified level
func NewWithLevel(level logrus.Level) *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)

	return &Logger{Logger: logger}
}
