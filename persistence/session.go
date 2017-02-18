package persistence

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/scoiatael/archai/types"
)

type Session interface {
	WriteEvent(string, []byte, map[string]string) error
	ReadEvents(string, string, int) ([]types.Event, error)
	Close()
}

type CassandraSession struct {
	*gocql.Session
}

const insertEvent = `INSERT INTO events (id, stream, blob, meta) VALUES (now(), ?, ?, ?)`

func (sess *CassandraSession) WriteEvent(stream string, blob []byte, meta map[string]string) error {
	return errors.Wrap(sess.Query(insertEvent, stream, blob, meta).Exec(), "Error writing event to Cassandra")
}

func eventFromRow(row map[string]interface{}) (types.Event, error) {
	var conv bool
	ev := types.Event{}
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

const readEvent = `SELECT id, blob, meta FROM events WHERE stream = ? AND id > ? LIMIT ?`

func (sess *CassandraSession) ReadEvents(stream string, cursor string, amount int) ([]types.Event, error) {
	iter := sess.Query(readEvent, stream, cursor, amount).Iter()
	rows, err := iter.SliceMap()
	events := make([]types.Event, len(rows))
	for i, r := range rows {
		events[i], err = eventFromRow(r)
		if err != nil {
			return events, errors.Wrap(err, "Conversion to Event failed")
		}
		events[i].Stream = stream
	}
	err = iter.Close()
	return events, errors.Wrap(err, "Failed readEvent")
}
