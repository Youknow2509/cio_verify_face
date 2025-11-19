package logger

import (
	"os"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger implements the ILogger interface using zap
type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// NewZapLogger creates a new zap logger instance
func NewZapLogger(config *domainConfig.LoggerConfig) (domainLogger.ILogger, *zap.Logger, error) {
	// Ensure log directory exists
	if err := os.MkdirAll(config.FolderStore, 0755); err != nil {
		return nil, nil, err
	}

	// Setup file output with rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.FolderStore + "/service.log",
		MaxSize:    config.FileMaxSize,    // megabytes
		MaxBackups: config.FileMaxBackups, // number of backups
		MaxAge:     config.FileMaxAge,     // days
		Compress:   config.Compress,
	})

	// Setup console output
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Create core
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			zapcore.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriter,
			zapcore.DebugLevel,
		),
	)

	// Build logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	return &ZapLogger{
		logger: logger,
		sugar:  sugar,
	}, logger, nil
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.sugar.Infow(msg, fields...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.sugar.Errorw(msg, fields...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	l.sugar.Warnw(msg, fields...)
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, fields ...interface{}) {
	l.sugar.Debugw(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	l.sugar.Fatalw(msg, fields...)
}
