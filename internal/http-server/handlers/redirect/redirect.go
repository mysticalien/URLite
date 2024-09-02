package redirect

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"

	resp "URLite/internal/lib/api/response"
	"URLite/internal/lib/logger/sl"
	"URLite/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter

// URLGetter — это интерфейс для получения URL по псевдониму.
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// New возвращает функцию-обработчик HTTP-запросов для перенаправления на оригинальный URL.
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		// Добавляем операцию и идентификатор запроса в логи для лучшей трассировки
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Извлекаем псевдоним из параметров URL
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		// Пытаемся получить оригинальный URL, используя псевдоним
		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("get url", slog.String("url", resURL))

		// Перенаправляем на полученный URL
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
