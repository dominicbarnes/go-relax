package relax

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is used for interacting with a CouchDB server.
type Client struct {
	http *http.Client
	base *url.URL
}

// Dial creates a client using the given addr URL as a base. Only scheme,
// hostname and authentication will be respected here.
func Dial(addr string) (*Client, error) {
	base, err := url.ParseRequestURI(addr)
	if err != nil {
		return nil, err
	}
	return &Client{new(http.Client), base}, nil
}

// Ping is a simple helper for checking that a server is reachable.
func (c *Client) Ping() error {
	res, err := c.request(http.MethodHead, c.resolve(nil, nil), nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return ErrInvalidResponse
	}
}

// Use returns a database-specific client.
func (c *Client) Use(db string) *DB {
	return &DB{c, db}
}

func (c *Client) resolve(path []string, query *url.Values) string {
	ref := url.URL{Path: "/" + strings.Join(path, "/")}
	if query != nil {
		ref.RawQuery = query.Encode()
	}
	return c.base.ResolveReference(&ref).String()
}

func (c *Client) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}

func (c *Client) decode(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	if v != nil {
		d := json.NewDecoder(res.Body)
		return d.Decode(v)
	}
	return nil
}

func (c *Client) error(res *http.Response) error {
	var cerr CouchDBError
	if err := c.decode(res, &cerr); err != nil {
		return err
	}
	return &cerr
}
