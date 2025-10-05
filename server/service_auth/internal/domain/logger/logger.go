package logger

import "errors"

/**
 * Interface logger service
 */
type ILogger interface {
	// Info logs
	Info(msg string, fields ...interface{})
	// Error logs
	Error(msg string, fields ...interface{})
	// Warn logs
	Warn(msg string, fields ...interface{})
	// Panic logs
	Panic(msg string, fields ...interface{})
	// Fatal logs
	Fatal(msg string, fields ...interface{})
}

/**
 * Manager instance logger
 */
var _vILogger ILogger

/**
 * Getter and setter instance logger
 */
func GetLogger() ILogger {
	return _vILogger
}

func SetLogger(logger ILogger) error {
	if logger == nil {
		return errors.New("instance set is nil")
	}
	if _vILogger != nil {
		return errors.New("instance already exists")
	}
	_vILogger = logger
	return nil
}