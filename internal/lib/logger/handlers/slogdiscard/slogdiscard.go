package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardLogger() *slog.Logger {
	return slog.New(&DiscardHandler{})
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle реализует метод интерфейса slog.Handler, но не делает ничего с логами.
func (h *DiscardHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

// WithAttrs возвращает новый обработчик с добавленными аттрибутами, но без изменений.
func (h *DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup возвращает новый обработчик с добавленной группой, но без изменений.
func (h *DiscardHandler) WithGroup(name string) slog.Handler {
	return h
}

// Всегда возвращает false, так как запись журнала игнорируется
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
