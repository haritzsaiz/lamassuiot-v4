package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime"

	"github.com/jakehl/goid"
)

const CtxSource = "REQ_SOURCE"
const CtxRequestID = "REQ_ID"

const CtxAuthMode = "REQ_AUTH_MODE"
const CtxAuthID = "REQ_AUTH_ID"

var LogFormatter = &Formatter{
	TimestampFormat: "2006-01-02 15:04:05",
	HideKeys:        true,
	FieldsOrder:     []string{"src", "auth-mode", "auth-id", "req-id", "service", "subsystem", "subsystem-provider"},
	CallerFirst:     true,
	CustomCallerFormatter: func(f *runtime.Frame) string {
		filename := path.Base(f.File)
		return fmt.Sprintf(" [%s %s():%d]", filename, f.Function, f.Line)
	},
	ReportCaller: false,
}

func SetupLogger(currentLevel Level, serviceID string, subsystem string) *Logger {
	var level Level

	var output io.Writer = os.Stdout
	if currentLevel == LevelNone {
		output = io.Discard
	}

	handler := NewFormatterHandler(output, level, LogFormatter)

	lSubsystem := New(handler).With(
		slog.String("service", serviceID),
		slog.String("subsystem", subsystem),
	)

	lSubsystem.Info(fmt.Sprintf("log level set to '%s'", level))
	return lSubsystem
}

func ConfigureLogger(ctx context.Context, logger *Logger) *Logger {
	logger = configureLoggerWitSourceAndCallerID(ctx, logger)
	logger = configureLoggerWithRequestID(ctx, logger)
	return logger
}

func configureLoggerWitSourceAndCallerID(ctx context.Context, log *Logger) *Logger {
	source := ""
	authMode := ""
	authID := ""

	sourceCtx := ctx.Value(string(CtxSource))
	if src, ok := sourceCtx.(string); ok {
		source = src
	}

	authIDCtx := ctx.Value(string(CtxAuthID))
	if id, ok := authIDCtx.(string); ok {
		authID = id
	}

	authModeCtx := ctx.Value(string(CtxAuthMode))
	if mode, ok := authModeCtx.(string); ok {
		authMode = mode
	}

	log = log.With("src", source)
	log = log.With("auth-mode", authMode)
	log = log.With("auth-id", authID)

	return log
}

func configureLoggerWithRequestID(ctx context.Context, log *Logger) *Logger {
	if !log.Enabled(context.Background(), LevelDebug) {
		return log
	}

	reqCtx := ctx.Value(CtxRequestID)
	if reqID, ok := reqCtx.(string); ok {
		return log.With("req-id", reqID)
	}

	return log.With("req-id", fmt.Sprintf("unset.%s", goid.NewV4UUID()))
}

func InitContext() context.Context {
	return context.WithValue(context.Background(), CtxRequestID, fmt.Sprintf("internal.%s", goid.NewV4UUID()))
}
