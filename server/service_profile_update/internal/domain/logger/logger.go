package logger

// ILogger defines the logging interface
type ILogger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	
	// Structured logging
	With(args ...interface{}) ILogger
}
