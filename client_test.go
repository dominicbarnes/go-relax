package relax_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/dominicbarnes/go-relax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const brokenURL = "http://localhost"

func TestDial(t *testing.T) {
	c, err := Dial(brokenURL)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestDialInvalidURL(t *testing.T) {
	c, err := Dial("")
	assert.EqualError(t, err, "parse : empty url")
	assert.Nil(t, c)
}

func TestClientPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodHead, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c, err := Dial(ts.URL)
	require.NoError(t, err)
	assert.NoError(t, c.Ping())
}

func TestClientPingNetworkError(t *testing.T) {
	c, err := Dial(brokenURL)
	require.NoError(t, err)
	assert.Error(t, c.Ping())
}

func TestClientPingServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send a response so we can test the decoding
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`boom`))
	}))
	defer ts.Close()

	c, err := Dial(ts.URL)
	require.NoError(t, err)
	assert.Error(t, c.Ping())
}
