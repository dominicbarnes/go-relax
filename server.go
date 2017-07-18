package relax

import (
	"net/http"
)

// ServerInfo contains CouchDB server information accessed via `GET /`.
type ServerInfo struct {
	CouchDB string                 `json:"couchDB,omitempty"`
	UUID    string                 `json:"uuid,omitempty"`
	Vendor  map[string]interface{} `json:"vendor,omitempty"`
	Version string                 `json:"version,omitempty"`
}

// Info requests basic server information.
func (c *Client) Info() (*ServerInfo, error) {
	var info ServerInfo
	err := c.relax(http.MethodGet, "/", nil, &info)
	return &info, err
}
