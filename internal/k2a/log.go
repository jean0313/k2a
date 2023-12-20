package k2a

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOG_FILE = "k2a.log"

// default info
var LOG_LEVEL = zap.NewAtomicLevel()

func InitLog() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, LOG_LEVEL),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), LOG_LEVEL),
	)
	zap.ReplaceGlobals(zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)))
}
