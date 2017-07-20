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
			request(t, r, http.MethodHead, "/test")

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

	t.Run("Server Error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		client, err := Dial(ts.URL)
		require.NoError(t, err)
		db := client.Use("test")

		exists, err := db.Exists()
		assert.Error(t, err)
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

func TestDBGet(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ensure we got the expected request
			request(t, r, http.MethodGet, "/db/doc")

			// return a document
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"_id":"document_1","_rev":"rev_1","answer":42}`))
		}))
		defer ts.Close()

		db := dial(t, ts).Use("db")

		type S struct {
			ID     string `json:"_id"`
			Rev    string `json:"_rev"`
			Answer int    `json:"answer"`
		}
		var s S

		assert.NoError(t, db.Get("doc", &s))
		assert.EqualValues(t, S{"document_1", "rev_1", 42}, s)
	})

	t.Run("Not Found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ensure we got the expected request
			request(t, r, http.MethodGet, "/db/doc")

			// return a document
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not_found","reason":"missing"}`))
		}))
		defer ts.Close()

		db := dial(t, ts).Use("db")

		type S struct {
			ID     string `json:"_id"`
			Rev    string `json:"_rev"`
			Answer int    `json:"answer"`
		}
		var s S

		assert.EqualError(t, db.Get("doc", &s), "[not_found] missing")
	})

	t.Run("Network Error", func(t *testing.T) {
		client, err := Dial(brokenURL)
		require.NoError(t, err)
		db := client.Use("test")

		assert.Error(t, db.Get("doc", nil))
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
