package relax_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dominicbarnes/go-relax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const integrationURL = "http://localhost:5984"
const brokenURL = "http://localhost"

func TestDial(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		c, err := Dial(integrationURL)
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		c, err := Dial("")
		assert.EqualError(t, err, "parse : empty url")
		assert.Nil(t, c)
	})
}

func TestClientPing(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ensure we got the expected request
			request(t, r, http.MethodHead, "/")

			// 200 OK => server is alive
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		c := dial(t, ts)
		assert.NoError(t, c.Ping())
	})

	t.Run("Server Error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(boom))
		defer ts.Close()

		c := dial(t, ts)
		assert.Error(t, c.Ping())
	})

	t.Run("Network Error", func(t *testing.T) {
		c, err := Dial(brokenURL)
		require.NoError(t, err)
		assert.Error(t, c.Ping())
	})
}

func TestClientIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	c, err := Dial(integrationURL)
	require.NoError(t, err)

	t.Run("Ping", func(t *testing.T) {
		assert.NoError(t, c.Ping())
	})
}

func dial(t *testing.T, ts *httptest.Server) *Client {
	c, err := Dial(ts.URL)
	require.NoError(t, err)
	return c
}

func request(t *testing.T, r *http.Request, method, path string) {
	assert.Equal(t, "application/json", r.Header.Get("Accept"))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, method, r.Method)
	assert.Equal(t, path, r.URL.Path)
}

func boom(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("boom!"))
}
