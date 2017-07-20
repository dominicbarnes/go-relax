package relax

// Document is an embeddable struct that conveniently adds id and rev to your
// own structs.
type Document struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}
