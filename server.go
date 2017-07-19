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
	res, err := c.request(http.MethodGet, c.resolve(nil, nil), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var info ServerInfo
		if err := c.decode(res, &info); err != nil {
			return nil, err
		}
		return &info, nil
	default:
		return nil, c.error(res)
	}
}
