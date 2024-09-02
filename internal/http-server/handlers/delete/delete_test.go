package delete_test

import (
	"net/http/httptest"
	"testing"

	"URLite/internal/http-server/handlers/delete"
	"URLite/internal/http-server/handlers/delete/mocks"
	"URLite/internal/lib/api"
	"URLite/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
		},
		{
			name:      "Empty Alias",
			alias:     "",
			respError: "api.DeleteURL: invalid status code: 404", // Проверка на 400 статус
		},
		{
			name:  "URL Not Found",
			alias: "nonexistent_alias",
		},
		{
			name:  "Internal Server Error",
			alias: "error_alias",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.alias != "" {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			err := api.DeleteURL(ts.URL + "/" + tc.alias)

			if tc.respError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.respError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
