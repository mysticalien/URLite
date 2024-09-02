package save_test

import (
	"URLite/internal/http-server/handlers/url/save"
	"URLite/internal/http-server/handlers/url/save/mocks"
	"URLite/internal/lib/logger/handlers/slogdiscard"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string // Имя теста
		alias     string // Псевдоним, который передается вместе с URL
		url       string // URL для сохранения
		respError string // Ожидаемое сообщение об ошибке в ответе
		mockError error  // Ошибка, которую должен вернуть мок-объект при попытке сохранения URL
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://go.dev/",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://go.dev/",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://go.dev/",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
		{
			name:      "Successful request without alias",
			alias:     "",
			url:       "https://example.com",
			respError: "",
		},
		{
			name:      "Alias with special characters",
			alias:     "special_!@#",
			url:       "https://example.com",
			respError: "",
		},
		{
			name:      "URL with query parameters",
			alias:     "query_params",
			url:       "https://example.com/search?q=test",
			respError: "",
		},
		{
			name:      "Duplicate URL Error",
			alias:     "duplicate_alias",
			url:       "https://example.com/",
			respError: "failed to add url",
			mockError: errors.New("duplicate URL error"),
		},
	}

	for _, tc := range cases {
		// Создаем локальную переменную tc, чтобы избежать проблем с параллельным выполнением тестов
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// Запускаем тест в параллельном режиме
			t.Parallel()

			// Создаем мок-объект для интерфейса URLSaver
			urlSaverMock := mocks.NewURLSaver(t)

			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomID := rng.Int63()

			// Настраиваем мок-объект в зависимости от тестового случая
			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(randomID, tc.mockError).
					Once()
			}

			// Создаем обработчик с использованием мок-объекта
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// Формируем входные данные для запроса
			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			// Создаем новый HTTP-запрос с использованием сформированных данных
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаем новый HTTP-респондер для записи ответа
			rr := httptest.NewRecorder()

			// Вызываем обработчик с созданным запросом
			handler.ServeHTTP(rr, req)

			// Проверяем, что код ответа соответствует ожидаемому (например, 200 OK)
			require.Equal(t, http.StatusOK, rr.Code)

			// Извлекаем тело ответа
			body := rr.Body.String()

			// Декодируем тело ответа в структуру
			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем, что в ответе содержится ожидаемое сообщение об ошибке
			require.Equal(t, tc.respError, resp.Error)

			// Проверяем, что поле Alias в ответе заполнено, если ожидается успешный результат
			if tc.respError == "" {
				require.NotEmpty(t, resp.Alias, "Alias should not be empty on success")
				if tc.alias != "" {
					require.Equal(t, tc.alias, resp.Alias, "Expected the alias in response to match the input alias")
				}
			}
		})
	}
}
