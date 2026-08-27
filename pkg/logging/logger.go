package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(env string) (*Logger, error) {
	logDir := "logs"
	logFile := filepath.Join(logDir, "app.log")

	// Create a directory for logs
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("Failed to create log directory: %w", err)
	}

	// Basic configuration of encoders
	var encoderConfig zapcore.EncoderConfig
	if env == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	// Выставляем понятный формат времени
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Setting up for the console
	consoleConfig := encoderConfig
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Setting for file
	fileConfig := encoderConfig
	fileConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var consoleEncoder, fileEncoder zapcore.Encoder

	if env == "production" {
		consoleEncoder = zapcore.NewJSONEncoder(consoleConfig)
		fileEncoder = zapcore.NewJSONEncoder(fileConfig)
	} else {
		consoleEncoder = zapcore.NewConsoleEncoder(consoleConfig)
		fileEncoder = zapcore.NewConsoleEncoder(fileConfig)
	}

	// Open the file for recording
	logFileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open log file: %w", err)
	}

	// Setting the logging level
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if env == "production" {
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Collecting kernels
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.Lock(logFileWriter), level)

	// Combining the output
	coreLogger := zapcore.NewTee(consoleCore, fileCore)

	// Creating a logger
	zapLogger := zap.New(coreLogger, zap.WithCaller(false))

	return &Logger{
		zapLogger.Sugar(),
	}, nil
}
