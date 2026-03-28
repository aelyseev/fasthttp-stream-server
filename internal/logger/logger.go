package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"
)

var (
	instance *Logger
	once     sync.Once
)

func Initialize(level slog.Level) {
	once.Do(func() {
		log.Println("Init logger")
		instance = newLogger(level)
	})
}

func GetLogger() *Logger {
	if instance == nil {
		panic("Logger is not initialized. Call logger.Initialize at first.")
	}
	return instance
}

type ctxKey struct{}

type Logger struct {
	l *slog.Logger
}

func newLogger(level slog.Level) *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return &Logger{
		l: slog.New(handler),
	}
}

func (lg *Logger) WithFields(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}

	existing := getAttrs(ctx)

	merged := make([]slog.Attr, 0, len(existing)+len(attrs))
	merged = append(merged, existing...)
	merged = append(merged, attrs...)

	return context.WithValue(ctx, ctxKey{}, merged)
}

func (lg *Logger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	lg.log(ctx, slog.LevelDebug, msg, attrs...)
}

func (lg *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	lg.log(ctx, slog.LevelInfo, msg, attrs...)
}

func (lg *Logger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	lg.log(ctx, slog.LevelWarn, msg, attrs...)
}

func (lg *Logger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	lg.log(ctx, slog.LevelError, msg, attrs...)
}

func (lg *Logger) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	ctxAttrs := getAttrs(ctx)

	all := make([]slog.Attr, 0, len(ctxAttrs)+len(attrs))
	all = append(all, ctxAttrs...)
	all = append(all, attrs...)

	lg.l.LogAttrs(ctx, level, msg, all...)
}

func getAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}

	v := ctx.Value(ctxKey{})
	if v == nil {
		return nil
	}

	attrs, ok := v.([]slog.Attr)
	if !ok {
		return nil
	}

	return attrs
}
