package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(logDir, appName, appVersion string) *zap.Logger {
	// Create log folder per app name
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}

	// Daily log file inside app-specific folder
	today := time.Now().Format("2006-01-02")
	logFile := logDir + "/" + today + ".log"

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

	baseLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Add app name/version to every log entry

	Logger = baseLogger.With(
		zap.String("app_name", appName),
		zap.String("app_version", appVersion),
	)

	return Logger
}
