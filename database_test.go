package relax_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dominicbarnes/go-relax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBExists(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ensure we got the expected request
			require.Equal(t, "application/json", r.Header.Get("Accept"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.Equal(t, http.MethodHead, r.Method)
			require.Equal(t, "/test", r.URL.Path)

			// 200 OK => db exists.
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		db := dial(t, ts).Use("test")

		exists, err := db.Exists()
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("False", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 404 Not Found => db does not exist.
			http.NotFound(w, r)
		}))
		defer ts.Close()

		db := dial(t, ts).Use("test")

		exists, err := db.Exists()
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Network Error", func(t *testing.T) {
		client, err := Dial(brokenURL)
		require.NoError(t, err)
		db := client.Use("test")

		exists, err := db.Exists()
		assert.Error(t, err)
		assert.False(t, exists)
	})
}

func TestDBIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	client, err := Dial(integrationURL)
	require.NoError(t, err)

	t.Run("Exists", func(t *testing.T) {
		db := client.Use("should_not_exist")

		exists, err := db.Exists()
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
