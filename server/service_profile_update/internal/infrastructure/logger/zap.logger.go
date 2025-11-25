package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/logger"
)

// ZapLogger implements ILogger using uber-go/zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger creates a new ZapLogger instance
func NewZapLogger(cfg *config.LoggerSetting) (domainLogger.ILogger, error) {
	// Parse log level - default to Info
	level := zapcore.InfoLevel

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Ensure log folder exists
	logFolder := cfg.FolderStore
	if logFolder == "" {
		logFolder = "./logs"
	}
	if err := os.MkdirAll(logFolder, 0755); err != nil {
		return nil, err
	}

	// Create file writer with rotation
	logFilePath := filepath.Join(logFolder, "profile_update.log")
	fileWriter := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    cfg.FileMaxSize,       // megabytes
		MaxBackups: cfg.FileMaxBackups,
		MaxAge:     cfg.FileMaxAge,        // days
		Compress:   cfg.Compress,
	}

	// Create core with both console and file output
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &ZapLogger{logger: logger.Sugar()}, nil
}

func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debugw(msg, args...)
}

func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warnw(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalw(msg, args...)
}

func (l *ZapLogger) With(args ...interface{}) domainLogger.ILogger {
	return &ZapLogger{logger: l.logger.With(args...)}
}
