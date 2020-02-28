package zlog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Field aliases zapcore.Field to limit the library scope into the code base
type Field = zapcore.Field

const (
	// LevelDebug has verbose message
	LevelDebug = "debug"
	// LevelInfo is default log level
	LevelInfo = "info"
	// LevelWarn is for logging messages about possible issues
	LevelWarn = "warn"
	// LevelError is for logging errors
	LevelError = "error"
	// LevelFatal is for logging fatal messages, process exits after
	LevelFatal = "fatal"
	// LevelPanic is for logging panic messages, it panics after
	LevelPanic = "panic"
)

var (
	// Int64 = zap.Int64
	Int64 = zap.Int64
	// Int32 = zap.Int32
	Int32 = zap.Int32
	// Int = zap.Int
	Int = zap.Int
	// Uint32 = zap.Uint32
	Uint32 = zap.Uint32
	// String = zap.String
	String = zap.String
	// Any = zap.Any
	Any = zap.Any
	// Err = zap.Error
	Err = zap.Error
	// NamedErr = zaNamedErrorp.
	NamedErr = zap.NamedError
	// Bool = zap.Bool
	Bool = zap.Bool
	// Duration = zaDurationp.
	Duration = zap.Duration
)

// LoggerConfig for the logger
type LoggerConfig struct {
	EnableConsole bool
	ConsoleJSON   bool
	ConsoleLevel  string
	EnableFile    bool
	FileJSON      bool
	FileLevel     string
	FileLocation  string
}

// Logger wraps zap Logger
type Logger struct {
	zap          *zap.Logger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoder(json bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if json {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// NewLogger creates the new Logger
func NewLogger(config *LoggerConfig) *Logger {
	cores := []zapcore.Core{}

	logger := &Logger{
		consoleLevel: zap.NewAtomicLevelAt(getZapLevel(config.ConsoleLevel)),
		fileLevel:    zap.NewAtomicLevelAt(getZapLevel(config.FileLevel)),
	}

	if config.EnableConsole {
		writer := zapcore.Lock(os.Stderr)
		core := zapcore.NewCore(getEncoder(config.ConsoleJSON), writer, logger.consoleLevel)
		cores = append(cores, core)
	}

	if config.EnableFile {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.FileLocation,
			MaxSize:  100,
			MaxAge:   28,
			Compress: true,
		})
		core := zapcore.NewCore(getEncoder(config.FileJSON), writer, logger.fileLevel)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger.zap = zap.New(combinedCore,
		zap.AddCaller(),
	)

	return logger
}

// Sugar wraps logger in a simpler interface
// it has worse performance (typed logger should be used mostly)
func (l *Logger) Sugar() *SugarLogger {
	return &SugarLogger{
		wrappedLogger: l,
		zapSugar:      l.zap.Sugar(),
	}
}

// ChangeLevels changes console and file log levels
func (l *Logger) ChangeLevels(config *LoggerConfig) {
	l.consoleLevel.SetLevel(getZapLevel(config.ConsoleLevel))
	l.fileLevel.SetLevel(getZapLevel(config.FileLevel))
}

// SetConsoleLevel sets the console log level
func (l *Logger) SetConsoleLevel(level string) {
	l.consoleLevel.SetLevel(getZapLevel(level))
}

// With wraps the logger with additional fields
func (l *Logger) With(fields ...Field) *Logger {
	newlogger := *l
	newlogger.zap = newlogger.zap.With(fields...)
	return &newlogger
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(message string, fields ...Field) {
	l.zap.Debug(message, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(message string, fields ...Field) {
	l.zap.Info(message, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(message string, fields ...Field) {
	l.zap.Warn(message, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(message string, fields ...Field) {
	l.zap.Error(message, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func (l *Logger) Fatal(message string, fields ...Field) {
	l.zap.Fatal(message, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panic(message string, fields ...Field) {
	l.zap.Panic(message, fields...)
}
