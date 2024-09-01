package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"

	"github.com/fatih/color"
	"log/slog"
)

// PrettyHandlerOptions содержит параметры для настройки PrettyHandler.
type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions // Опции для JSON-обработчика
}

// PrettyHandler реализует обработку логов с цветным выводом в консоль.
type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger // Стандартный логгер для вывода в консоль
	attrs []slog.Attr    // Дополнительные атрибуты для логирования
}

// NewPrettyHandler создает новый PrettyHandler с заданным io.Writer для вывода.
func (opts PrettyHandlerOptions) NewPrettyHandler(
	out io.Writer,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

// Handle обрабатывает запись лога и выводит ее в консоль с цветным форматированием.
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	// Форматирование уровня логирования
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Сбор атрибутов
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	// Форматирование атрибутов в JSON
	var b []byte
	var err error
	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	// Форматирование и вывод сообщения
	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)
	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

// WithAttrs добавляет дополнительные атрибуты к логам и возвращает новый обработчик.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

// WithGroup добавляет групповые атрибуты к логам.
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// TODO: implement
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
