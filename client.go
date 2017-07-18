package relax

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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
	return c.relax(http.MethodHead, "/", nil, nil)
}

// Use returns a database-specific client.
func (c *Client) Use(db string) *DB {
	return &DB{c, db}
}

func (c *Client) relax(method, ref string, body io.Reader, v interface{}) error {
	req, err := c.request(method, ref, body)
	if err != nil {
		return err
	}
	return c.response(req, v)
}

func (c *Client) url(input string) (string, error) {
	ref, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	return c.base.ResolveReference(ref).String(), nil
}

func (c *Client) request(method, ref string, body io.Reader) (*http.Request, error) {
	url, err := c.url(ref)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")

	return req, nil
}

func (c *Client) response(req *http.Request, v interface{}) error {
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if v != nil {
		d := json.NewDecoder(res.Body)
		return d.Decode(v)
	}
	return nil
}
