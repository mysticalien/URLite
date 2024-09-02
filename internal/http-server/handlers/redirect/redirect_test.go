package redirect_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"URLite/internal/http-server/handlers/redirect"
	"URLite/internal/http-server/handlers/redirect/mocks"
	"URLite/internal/lib/api"
	"URLite/internal/lib/logger/handlers/slogdiscard"
	"URLite/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com",
		},
		{
			name:      "Empty Alias",
			alias:     "",
			respError: "api.GetRedirect: invalid status code: 404",
		},
		{
			name:      "URL Not Found",
			alias:     "nonexistent_alias",
			respError: "api.GetRedirect: invalid status code: 200",
			mockError: storage.ErrURLNotFound,
		},
		{
			name:      "Internal Server Error",
			alias:     "error_alias",
			respError: "api.GetRedirect: invalid status code: 200",
			mockError: errors.New("internal error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.alias != "" {
				urlGetterMock.On("GetURL", tc.alias).Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			if tc.respError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.respError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.url, redirectedToURL)
			}
		})
	}
}
