package logger

import (
	"os"

	"github.com/rezatg/payment-gateway/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents a configured Zap logger
var Logger *zap.Logger

// SugaredLogger provides a sugared interface for easier logging
var SugaredLogger *zap.SugaredLogger

// Init initializes the global Zap logger
func InitLogger() error {
	logLevel := config.GetEnv("LOG_LEVEL", "info")
	env := config.GetEnv("APP_ENV", "development")

	// Configure core
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure encoder for JSON output
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if env == "production" {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)

	// Initialize logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	SugaredLogger = Logger.Sugar()

	return nil
}

// Debug logs a debug message with optional fields
func Debug(message string, fields ...interface{}) {
	SugaredLogger.Debugw(message, fields...)
}

// Info logs an info message with optional fields
func Info(message string, fields ...interface{}) {
	SugaredLogger.Infow(message, fields...)
}

// Warn logs a warning message with optional fields
func Warn(message string, fields ...interface{}) {
	SugaredLogger.Warnw(message, fields...)
}

// Error logs an error message with optional fields
func Error(message string, err error, fields ...interface{}) {
	fields = append(fields, "error", err)
	SugaredLogger.Errorw(message, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(message string, err error, fields ...interface{}) {
	fields = append(fields, "error", err)
	SugaredLogger.Fatalw(message, fields...)
}
