package relax

import (
	"errors"
	"fmt"
)

// ErrInvalidResponse is returned when something unexpected happens when
// interacting with CouchDB.
var ErrInvalidResponse = errors.New("invalid response")

// CouchDBError is used to decode errors reported by CouchDB.
type CouchDBError struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

// Error implements the error interface by returning the code and reason as a
// string.
func (e CouchDBError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Reason)
}
