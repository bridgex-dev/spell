package spell

import (
	"io"
	"log"
	"os"
)

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = iota
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

type Logger interface {
	Logf(level LogLevel, msg string, args ...interface{})
	SetLevel(level LogLevel)
}

type DefaultLogger struct {
	logger *log.Logger
	level  LogLevel
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "spell:", log.LstdFlags),
		level:  InfoLevel,
	}
}

func NewDefaultLoggerWithWriter(w io.Writer) *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(w, "spell:", log.LstdFlags),
		level:  InfoLevel,
	}
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *DefaultLogger) Logf(level LogLevel, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg = "[" + level.String() + "] " + msg
	l.logger.Printf(msg, args...)
}
