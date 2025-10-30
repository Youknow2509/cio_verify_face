package logger

import (
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/**
 * Struct logger implementation
 */
type ZapLogger struct {
	logger *zap.Logger
}

/**
 * Initialize a new logger service
 */
func NewZapLogger(initializer *ZapLoggerInitializer) (domainLogger.ILogger, error) {
	// only log warnings and errors
	levelLoggerWarn := zapcore.WarnLevel
	levelLoggerError := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	// define lumberjack
	lumberjackWarning := lumberjack.Logger{
		Filename:   initializer.FolderStore + "/warning.log",
		MaxSize:    initializer.FileMaxSize,    // megabytes
		MaxBackups: initializer.FileMaxBackups, // number of backups
		MaxAge:     initializer.FileMaxAge,     // days
		Compress:   initializer.Compress,       // compress the backups
	}
	lumberjackError := lumberjack.Logger{
		Filename:   initializer.FolderStore + "/error.log",
		MaxSize:    initializer.FileMaxSize,    // megabytes
		MaxBackups: initializer.FileMaxBackups, // number of backups
		MaxAge:     initializer.FileMaxAge,     // days
		Compress:   initializer.Compress,       // compress the backups
	}
	lumberjackDebug := lumberjack.Logger{
		Filename:   initializer.FolderStore + "/debug.log",
		MaxSize:    initializer.FileMaxSize,    // megabytes
		MaxBackups: initializer.FileMaxBackups, // number of backups
		MaxAge:     initializer.FileMaxAge,     // days
		Compress:   initializer.Compress,       // compress the backups
	}
	// create zap core
	core := zapcore.NewTee(
		// debug core
		zapcore.NewCore(getJSONEncoder(), zapcore.AddSync(&lumberjackDebug), zapcore.DebugLevel),
		// warning core
		zapcore.NewCore(getJSONEncoder(), zapcore.AddSync(&lumberjackWarning), levelLoggerWarn),
		// error core
		zapcore.NewCore(getJSONEncoder(), zapcore.AddSync(&lumberjackError), levelLoggerError),
	)
	// create logger with core
	logger := zap.New(core).WithOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return &ZapLogger{
		logger: logger,
	}, nil
}

// Error implements logger.ILogger.
func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	zapFields := convertFields(fields...)
	l.logger.Error(msg, zapFields...)
}

// Fatal implements logger.ILogger.
func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	zapFields := convertFields(fields...)
	l.logger.Fatal(msg, zapFields...)
}

// Info implements logger.ILogger.
func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	zapFields := convertFields(fields...)
	l.logger.Info(msg, zapFields...)
}

// Panic implements logger.ILogger.
func (l *ZapLogger) Panic(msg string, fields ...interface{}) {
	zapFields := convertFields(fields...)
	l.logger.Panic(msg, zapFields...)
}

// Warn implements logger.ILogger.
func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	zapFields := convertFields(fields...)
	l.logger.Warn(msg, zapFields...)
}

/**
 * Struct data initialization logger service
 */
type ZapLoggerInitializer struct {
	FolderStore    string `json:"FolderStore" yaml:"folder_store"`
	FileMaxSize    int    `json:"FileMaxSize" yaml:"file_max_size"`
	FileMaxBackups int    `json:"FileMaxBackups" yaml:"file_max_backups"`
	FileMaxAge     int    `json:"FileMaxAge" yaml:"file_max_age"`
	Compress       bool   `json:"Compress" yaml:"compress"`
}

// ===================================================
// 					Helper
// ===================================================

/**
 * Helper convert fields to zap fields
 */
func convertFields(fields ...interface{}) []zap.Field {
	if len(fields) == 0 {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	// parse fields interface to zap fields
	for _, field := range fields {
		switch v := field.(type) {
		case string:
			zapFields = append(zapFields, zap.String("field", v))
		case int:
			zapFields = append(zapFields, zap.Int("field", v))
		case int64:
			zapFields = append(zapFields, zap.Int64("field", v))
		case bool:
			zapFields = append(zapFields, zap.Bool("field", v))
		case float64:
			zapFields = append(zapFields, zap.Float64("field", v))
		case error:
			zapFields = append(zapFields, zap.Error(v))
		default:
			zapFields = append(zapFields, zap.Any("field", v))
		}
	}
	return zapFields
}

/**
 * Helper get json encoder
 */
func getJSONEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(getEncoderConfig())
}

/**
 * Helper get encoder config
 */
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// ===== Keys =====
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		// ===== Values =====
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		// ===== Options =====
		// v.v
	}
}
