package relax_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dominicbarnes/go-relax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientInfo(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ensure we got the expected request
			require.Equal(t, "application/json", r.Header.Get("Accept"))
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/", r.URL.Path)

			// send a response so we can test the decoding
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"couchdb":"Welcome","uuid":"85fb71bf700c17267fef77535820e371","vendor":{"name":"The Apache Software Foundation","version":"1.3.1"},"version":"1.3.1"}`))
		}))
		defer ts.Close()

		c := dial(t, ts)

		info, err := c.Info()
		assert.NoError(t, err)
		assert.EqualValues(t, &ServerInfo{
			CouchDB: "Welcome",
			UUID:    "85fb71bf700c17267fef77535820e371",
			Vendor: map[string]interface{}{
				"name":    "The Apache Software Foundation",
				"version": "1.3.1",
			},
			Version: "1.3.1",
		}, info)
	})

	t.Run("Server Error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// send a response so we can test the decoding
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`boom`))
		}))
		defer ts.Close()

		c := dial(t, ts)

		info, err := c.Info()
		assert.Error(t, err)
		assert.Nil(t, info)
	})

	t.Run("Network Error", func(t *testing.T) {
		c, err := Dial(brokenURL)
		require.NoError(t, err)

		info, err := c.Info()
		assert.Error(t, err)
		assert.Nil(t, info)
	})
}

func TestServerIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	c, err := Dial(integrationURL)
	require.NoError(t, err)

	t.Run("Info", func(t *testing.T) {
		info, err := c.Info()
		assert.NoError(t, err)
		assert.EqualValues(t, &ServerInfo{
			CouchDB: "Welcome",
			Vendor:  map[string]interface{}{"name": "The Apache Software Foundation"},
			Version: "2.0.0",
		}, info)
	})
}
