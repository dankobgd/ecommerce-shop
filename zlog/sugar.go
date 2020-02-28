package zlog

import "go.uber.org/zap"

// SugarLogger holds zap sugar logger for the ease of use
type SugarLogger struct {
	wrappedLogger *Logger
	zapSugar      *zap.SugaredLogger
}

// Debug logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *SugarLogger) Debug(msg string, keyValuePairs ...interface{}) {
	l.zapSugar.Debugw(msg, keyValuePairs...)
}

// Info logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *SugarLogger) Info(msg string, keyValuePairs ...interface{}) {
	l.zapSugar.Infow(msg, keyValuePairs...)
}

// Error logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *SugarLogger) Error(msg string, keyValuePairs ...interface{}) {
	l.zapSugar.Errorw(msg, keyValuePairs...)
}

// Warn logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (l *SugarLogger) Warn(msg string, keyValuePairs ...interface{}) {
	l.zapSugar.Warnw(msg, keyValuePairs...)
}
