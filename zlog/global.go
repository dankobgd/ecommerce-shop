package zlog

import (
	"go.uber.org/zap"
)

var globalLogger *Logger

// LogFunc is default log signature
type LogFunc func(string, ...Field)

// Debug LogFunc = defaultDebugLog
var Debug LogFunc = defaultDebugLog

// Info LogFunc = defaultInfoLog
var Info LogFunc = defaultInfoLog

// Warn LogFunc = defaultWarnLog
var Warn LogFunc = defaultWarnLog

// Error LogFunc = defaultErrorLog
var Error LogFunc = defaultErrorLog

// Fatal LogFunc = defaultErrorLog
var Fatal LogFunc = defaultErrorLog

// Panic LogFunc = defaultPanicLog
var Panic LogFunc = defaultPanicLog

// InitGlobalLogger inits the logger for global use
func InitGlobalLogger(logger *Logger) {
	glob := *logger
	glob.zap = glob.zap.WithOptions(zap.AddCallerSkip(1))
	globalLogger = &glob
	Debug = globalLogger.Debug
	Info = globalLogger.Info
	Warn = globalLogger.Warn
	Error = globalLogger.Error
	Fatal = globalLogger.Fatal
	Panic = globalLogger.Panic
}
