package relax_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dominicbarnes/go-relax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBExistsTrue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure we got the expected request
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, http.MethodHead, r.Method)
		require.Equal(t, "/test", r.URL.Path)

		// send back a 200 OK to indicate the db exists.
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, err := Dial(ts.URL)
	require.NoError(t, err)
	db := client.Use("test")

	exists, err := db.Exists()
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDBExistsFalse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure we got the expected request
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.Equal(t, http.MethodHead, r.Method)
		require.Equal(t, "/test", r.URL.Path)

		// send back a 404 Not Found to indicate the db does not exist.
		http.NotFound(w, r)
	}))

	client, err := Dial(ts.URL)
	require.NoError(t, err)
	db := client.Use("test")

	exists, err := db.Exists()
	assert.NoError(t, err)
	assert.False(t, exists)
}
