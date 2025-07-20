package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type Handler func(msg string, args ...any)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logger struct {
	level    slog.Level
	slogger  *slog.Logger
	handlers map[slog.Level][]Handler
	mu       sync.RWMutex
}

var l *logger

func New() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	l = &logger{
		slogger:  slog.New(handler),
		handlers: make(map[slog.Level][]Handler),
	}
}

func AddLogger(logger Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.handlers[slog.LevelDebug] = append(l.handlers[slog.LevelDebug], logger.Debug)
	l.handlers[slog.LevelInfo] = append(l.handlers[slog.LevelInfo], logger.Info)
	l.handlers[slog.LevelWarn] = append(l.handlers[slog.LevelWarn], logger.Warn)
	l.handlers[slog.LevelError] = append(l.handlers[slog.LevelError], logger.Error)
}

func SetLevel(level slog.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = level
	l.slogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

func Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

func Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

func Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, args...)
}

func (l *logger) log(level slog.Level, msg string, args ...any) {
	if level >= l.level {
		switch level {
		case slog.LevelDebug:
			l.slogger.Debug(msg, args...)
		case slog.LevelInfo:
			l.slogger.Info(msg, args...)
		case slog.LevelWarn:
			l.slogger.Warn(msg, args...)
		case slog.LevelError:
			l.slogger.Error(msg, args...)
		}
	}

	l.callHandlers(level, msg, args...)
}

func (l *logger) callHandlers(level slog.Level, msg string, args ...any) {
	l.mu.RLock()
	handlers, exists := l.handlers[level]
	l.mu.RUnlock()

	if !exists {
		return
	}

	for _, handler := range handlers {
		go func(h Handler) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("Recovered from panic in log handler", "level", level, "message", msg, "args", args, "error", r)
				}
			}()
			h(msg, args...)
		}(handler)
	}
}

func LevelFromString(levelStr string) (slog.Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return -1, fmt.Errorf("Unknown log level: '%s'", levelStr)
	}
}
