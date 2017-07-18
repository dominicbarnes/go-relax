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
	req, err := db.client.request(http.MethodHead, db.path(), nil)
	if err != nil {
		return false, err
	}
	res, err := db.client.http.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK, nil
}

func (db *DB) path() string {
	return "/" + db.name
}
