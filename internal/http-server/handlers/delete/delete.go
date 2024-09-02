package delete

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

// URLDeleter — это интерфейс для удаления URL по псевдониму.
type URLDeleter interface {
	DeleteURL(alias string) error
}

// New возвращает функцию-обработчик HTTP-запросов для удаления URL по псевдониму.
func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		// Добавляем операцию и идентификатор запроса в логи для лучшей трассировки
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Извлекаем псевдоним из параметров URL
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("empty alias")
			render.JSON(w, r, resp.Error("incorrect request"))
			return
		}

		// Пытаемся удалить URL по псевдониму
		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete URL", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("URL successfully deleted", slog.String("alias", alias))

		// Возвращаем успешный ответ
		render.JSON(w, r, resp.OK())
	}
}
