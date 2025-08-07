package logger

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/natefinch/lumberjack"
	"github.com/ntdat104/go-finance-dataset/pkg/datetime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger atomic.Value

func InitDefault() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap.NewProduction has error: %v", err)
	}
	setGlobalLogger(logger)
}

func InitProduction(filePath string) {
	// Create log folder per app name
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0755)
	}

	// Daily log file inside app-specific folder
	today := datetime.ConvertCurrentLocalTimeToString(datetime.YYYY_MM_DD)
	logFile := filePath + "/" + today + ".log"

	// File writer with rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // MB
		MaxBackups: 7,
		MaxAge:     30,   // days
		Compress:   true, // .gz
	})

	// Console writer
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Encoder config
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	// Combine file and console
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriter, zapcore.InfoLevel),
		zapcore.NewCore(jsonEncoder, consoleWriter, zapcore.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	setGlobalLogger(logger)
}

// Sync flushs any buffered log entries. It should be call before program exit.
func Sync() error {
	return getGlobalLogger().Sync()
}

// Debug logs a message at Debug level.
func Debug(msg string, fields ...zap.Field) {
	getGlobalLogger().Debug(msg, fields...)
}

// Info logs a message at Info level.
func Info(msg string, fields ...zap.Field) {
	getGlobalLogger().Info(msg, fields...)
}

// Error logs a message at Error level.
func Error(msg string, fields ...zap.Field) {
	getGlobalLogger().Error(msg, fields...)
}

// Warn logs a message at Warn level.
func Warn(msg string, fields ...zap.Field) {
	getGlobalLogger().Warn(msg, fields...)
}

// Fatal logs a message at Fatal level.
func Fatal(msg string, fields ...zap.Field) {
	getGlobalLogger().Fatal(msg, fields...)
}

func setGlobalLogger(logger *zap.Logger) {
	globalLogger.Store(logger)
}

func getGlobalLogger() *zap.Logger {
	return globalLogger.Load().(*zap.Logger)
}
