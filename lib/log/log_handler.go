package log

import (
	"context"
	"log/slog"

	"github.com/kujilabo/redstart/lib/domain"
)

type LogHandler struct {
	slog.Handler
}

var (
	LoggerNameKey = "logger_name"
)

const (
	LoggerNameContextKey domain.ContextKey = "LoggerNameContextKey"
)

func (h *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	loggerName, ok := ctx.Value(LoggerNameContextKey).(string)
	if ok {
		record.AddAttrs(slog.String(LoggerNameKey, loggerName))
	}

	return h.Handler.Handle(ctx, record)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		Handler: h.Handler.WithGroup(name),
	}
}
