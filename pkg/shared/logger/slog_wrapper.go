package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
)

type Level int

func (l Level) Level() Level { return Level(l) }
func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l *Level) Parse(s string) {
	switch strings.ToLower(s) {
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warn":
		*l = LevelWarn
	case "error":
		*l = LevelError
	case "trace":
		*l = LevelTrace
	case "fatal":
		*l = LevelFatal
	case "none":
		*l = LevelNone
	default:
		*l = LevelInfo
	}
}

func (l *Level) MarshalText() ([]byte, error) {
	if l == nil {
		return []byte("unknown"), nil
	}
	return []byte(l.String()), nil
}

type Leveler interface {
	Level() Level
}

type LevelVar struct {
	val atomic.Int64
}

func (v *LevelVar) Level() Level {
	return Level(v.val.Load())
}

func (v *LevelVar) Set(l Level) {
	v.val.Store(int64(l))
}

const (
	LevelTrace Level = -8 // Trace: more verbose than Debug
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12 // Fatal: used for critical errors that cause the program to exit
	LevelNone  Level = 16
)

type Logger struct {
	slog *slog.Logger
}

// Factory for a new instance
func New(handler slog.Handler) *Logger {
	return &Logger{slog: slog.New(handler)}
}

func SetDefault(log *Logger) {
	slog.SetDefault(log.slog)
}

// ====== Log methods ======
func (l *Logger) Info(msg string, args ...any)  { l.slog.Info(msg, args...) }
func (l *Logger) Debug(msg string, args ...any) { l.slog.Debug(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.slog.Warn(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.slog.Error(msg, args...) }
func (l *Logger) Trace(msg string, args ...any) {
	l.slog.Log(context.Background(), slog.Level(LevelTrace), msg, args...)
}
func (l *Logger) Fatal(msg string, args ...any) {
	l.slog.Log(context.Background(), slog.Level(LevelFatal), msg, args...)
	os.Exit(1) // Exit after logging fatal error
}

func (l *Logger) Infof(format string, args ...any)  { l.slog.Info(fmt.Sprintf(format, args...)) }
func (l *Logger) Debugf(format string, args ...any) { l.slog.Debug(fmt.Sprintf(format, args...)) }
func (l *Logger) Warnf(format string, args ...any)  { l.slog.Warn(fmt.Sprintf(format, args...)) }
func (l *Logger) Errorf(format string, args ...any) { l.slog.Error(fmt.Sprintf(format, args...)) }
func (l *Logger) Tracef(format string, args ...any) {
	l.slog.Log(context.Background(), slog.Level(LevelTrace), fmt.Sprintf(format, args...))
}
func (l *Logger) Fatalf(format string, args ...any) {
	l.slog.Log(context.Background(), slog.Level(LevelFatal), fmt.Sprintf(format, args...))
	os.Exit(1) // Exit after logging fatal error
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.slog.With(args...)}
}
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	return l.slog.Enabled(ctx, slog.Level(level))
}

func Info(msg string, args ...any)  { slog.Info(msg, args...) }
func Debug(msg string, args ...any) { slog.Debug(msg, args...) }
func Warn(msg string, args ...any)  { slog.Warn(msg, args...) }
func Error(msg string, args ...any) { slog.Error(msg, args...) }
func Trace(msg string, args ...any) {
	slog.Log(context.Background(), slog.Level(LevelTrace), msg, args...)
}

func Fatal(msg string, args ...any) {
	slog.Log(context.Background(), slog.Level(LevelFatal), msg, args...)
	os.Exit(1) // Exit after logging fatal error
}

func Infof(format string, args ...any)  { slog.Info(fmt.Sprintf(format, args...)) }
func Debugf(format string, args ...any) { slog.Debug(fmt.Sprintf(format, args...)) }
func Warnf(format string, args ...any)  { slog.Warn(fmt.Sprintf(format, args...)) }
func Errorf(format string, args ...any) { slog.Error(fmt.Sprintf(format, args...)) }
func Tracef(format string, args ...any) {
	slog.Log(context.Background(), slog.Level(LevelTrace), fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...any) {
	slog.Log(context.Background(), slog.Level(LevelFatal), fmt.Sprintf(format, args...))
	os.Exit(1) // Exit after logging fatal error
}

func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "error":
		return Level(LevelError), nil
	case "warn", "warning":
		return Level(LevelWarn), nil
	case "info":
		return Level(LevelInfo), nil
	case "debug":
		return Level(LevelDebug), nil
	case "trace":
		return Level(LevelTrace), nil
	case "fatal":
		return Level(LevelFatal), nil
	}

	var l Level
	return l, fmt.Errorf("not a valid slog Level: %q", lvl)
}
