package relax

import (
	"net/http"
)

// DB is used for interacting with a specific CouchDB database.
type DB struct {
	client *Client
	name   string
}

// Exists checks for the database's existence via `HEAD /:db`.
func (db *DB) Exists() (bool, error) {
	url := db.client.resolve([]string{db.name}, nil)
	_, res, err := db.client.request(http.MethodHead, url, nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, ErrInvalidResponse
	}
}
