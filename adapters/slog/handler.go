package slog

import (
	"context"
	"log/slog"

	"github.com/cnmax/gologging-ext/core"
)

type Handler struct {
	writer core.Writer
	level  slog.Level
}

func NewHandler(writer core.Writer, level slog.Level) *Handler {
	return &Handler{
		writer: writer,
		level:  level,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	entry := &core.Entry{
		Time:    r.Time,
		Level:   r.Level.String(),
		Message: r.Message,
		Fields:  map[string]interface{}{},
	}

	source := r.Source()

	if source != nil {
		entry.Function = source.Function
		entry.File = source.File
		entry.Line = source.Line
	}

	r.Attrs(func(attr slog.Attr) bool {
		entry.Fields[attr.Key] = attr.Value.Any()
		return true
	})

	return h.writer.Write(entry)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}
