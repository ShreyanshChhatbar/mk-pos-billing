package logger

import (
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var consoleLoggingEnabled = true

// ConsoleLoggingEnabled returns whether zap console output is allowed.
func ConsoleLoggingEnabled() bool {
	return consoleLoggingEnabled
}

// Init initializes Zap Logger + Sentry for error reporting.
func Init(env string, sentryDSN string) {
	envNormalized := strings.ToLower(strings.TrimSpace(env))
	consoleLoggingEnabled = envNormalized == "development"

	initSentry(env, sentryDSN)

	// File-based Rotating Logger
	writer, _ := rotatelogs.New(
		"logs/app-%Y%m%d.log",
		rotatelogs.WithLinkName("logs/current.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(writer), zap.DebugLevel)

	cores := []zapcore.Core{fileCore} // ✅ Always log to file

	if consoleLoggingEnabled {
		// ✅ Also log to console in development
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel)
		cores = append(cores, consoleCore)
	}

	// Sentry Core - only for ERROR level and above
	if sentryDSN != "" {
		sentryCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(sentryWriter{}),
			zapcore.ErrorLevel, // Only Errors and above go to Sentry
		)
		cores = append(cores, sentryCore)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	zap.ReplaceGlobals(logger)

	zap.L().Info("Logger initialized", zap.String("env", env))
}

func initSentry(env, sentryDSN string) {
	if sentryDSN == "" {
		return
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDSN,
		Environment:      env,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		TracesSampler: func(ctx sentry.SamplingContext) float64 {
			if ctx.Span != nil && ctx.Span.Name == "GET /health" {
				return 0
			}
			return 1.0
		},
	})
	if err != nil {
		panic("sentry.Init failed: " + err.Error())
	}
}

// sentryWriter sends logs to Sentry as events
type sentryWriter struct{}

func (sentryWriter) Write(p []byte) (n int, err error) {
	sentry.CurrentHub().CaptureMessage(string(p))
	return len(p), nil
}

// Writer redirects panic stack traces to Zap logger & Sentry
type Writer struct{}

func (w Writer) Write(p []byte) (n int, err error) {
	zap.L().Error("PANIC", zap.ByteString("stacktrace", p))
	sentry.CaptureMessage(string(p))
	return len(p), nil
}

// RecoverAndReport recovers from panic, logs it and sends to Sentry
func RecoverAndReport() {
	if err := recover(); err != nil {
		zap.L().Error("Recovered from panic", zap.Any("error", err))
		sentry.CurrentHub().Recover(err)
		sentry.Flush(2 * time.Second)
	}
}

// Sync flushes Zap and Sentry
func Sync() {
	_ = zap.L().Sync()
	sentry.Flush(2 * time.Second)
}
