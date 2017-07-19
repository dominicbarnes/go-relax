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
	res, err := db.client.request(http.MethodHead, url, nil)
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

func (db *DB) Get(id string, v interface{}) error {
	url := db.client.resolve([]string{db.name, id}, nil)
	res, err := db.client.request(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return db.client.decode(res, v)
	default:
		return db.client.error(res)
	}
}
