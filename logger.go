package sllogger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// Logger is a struct that wraps around a zerolog.Logger instance.
type Logger struct {
	*zerolog.Logger
}

// RequestID is a type used to represent a unique identifier for a request.
type RequestID struct{}

const (
	requestID = "request_id"
)

var (
	defaultLogger *zerolog.Logger
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	debug := os.Getenv("DEBUG")
	var isDebug bool

	if debug != "" {
		if isDebug, err = strconv.ParseBool(debug); err != nil {
			fmt.Printf("error while parsing env 'debug' flag: %s\n", err.Error())
		}
	}

	if isDebug {
		logLevel = zerolog.DebugLevel.String()
	}

	var zeroLogger zerolog.Logger
	var writer io.Writer = os.Stdout

	if isDebug {
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	zeroLogger = zerolog.New(writer)

	zeroLogger = zeroLogger.
		With().
		Caller().
		Timestamp().
		Str("goversion", runtime.Version()).
		Str("host", hostname).
		Logger()

	appName := os.Getenv("APP_NAME")
	if appName != "" {
		zeroLogger = zeroLogger.With().Str("app", appName).Logger()
	}

	// Set proper loglevel based on config
	switch strings.ToLower(logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // log info and above by default
	}

	defaultLogger = &zeroLogger
}

// fromCtx is a function for getting logger from context or defaultLogger.
func fromCtx(ctx context.Context) *zerolog.Logger {
	logger := defaultLogger

	if ctx == nil {
		return logger
	}

	reqID := ctx.Value(RequestID{}).(string)
	if reqID != "" {
		*logger = logger.With().Str(requestID, reqID).Logger()
	}

	return logger
}

// Info is a function which return log event with level info.
func Info(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Info()
}

// Warn is a function which return log event with level warn.
func Warn(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Warn()
}

// Error is a function which return log event with level error.
func Error(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Error()
}

// Err is a function which return log event with level error.
func Err(ctx context.Context, err error) *zerolog.Event {
	return fromCtx(ctx).Err(err)
}

// Fatal is a function which return log event with level fatal.
func Fatal(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Fatal()
}

// Debug is a function which return log event with level debug.
func Debug(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Debug()
}

// Panic is a function which return log event with level panic.
func Panic(ctx context.Context) *zerolog.Event {
	return fromCtx(ctx).Panic()
}

// Get is a function for getting defaultLogger.
func Get() *Logger {
	return &Logger{Logger: defaultLogger}
}
