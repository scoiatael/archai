package persistence

import (
	"fmt"

	"github.com/gocql/gocql"
)

type Event struct {
	ID     string
	Stream string
	Blob   []byte
	Meta   map[string]string
}

func eventFromRow(row map[string]interface{}) (Event, error) {
	var conv bool
	ev := Event{}
	id, conv := row["id"].(gocql.UUID)
	if !conv {
		return ev, fmt.Errorf("Failed converting %v to UUID", row["id"])
	}
	ev.ID = id.String()
	ev.Blob, conv = row["blob"].([]byte)
	if !conv {
		return ev, fmt.Errorf("Failed converting %v to blob", row["blob"])
	}
	ev.Meta, conv = row["meta"].(map[string]string)
	if !conv {
		return ev, fmt.Errorf("Failed converting %v to map", row["map"])
	}
	return ev, nil
}
