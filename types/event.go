package types

import (
	"github.com/scoiatael/archai/simplejson"
)

type Event struct {
	ID     string
	Stream string
	Blob   []byte
	Meta   map[string]string
}

func EventToJson(e Event) ([]byte, error) {
	object := make(simplejson.Object)
	object["id"] = e.ID
	object["stream"] = e.Stream
	payload, err := simplejson.Read(e.Blob)
	if err != nil {
		return []byte{}, err
	}
	object["payload"] = payload
	return simplejson.Write(object)

}
