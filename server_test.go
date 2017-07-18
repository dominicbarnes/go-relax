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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"couchdb":"Welcome","uuid":"85fb71bf700c17267fef77535820e371","vendor":{"name":"The Apache Software Foundation","version":"1.3.1"},"version":"1.3.1"}`))
	}))
	c, err := Dial(ts.URL)
	require.NoError(t, err)
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
}
