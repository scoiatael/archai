package persistence

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

type Session interface {
	WriteEvent(string, []byte, map[string]string) error
	ReadEvents(string, string, int) ([]Event, error)
}

type CassandraSession struct {
	session *gocql.Session
}

const insertEvent = `INSERT INTO events (id, stream, blob, meta) VALUES (now(), ?, ?, ?)`

func (sess *CassandraSession) WriteEvent(stream string, blob []byte, meta map[string]string) error {
	return errors.Wrap(sess.session.Query(insertEvent, stream, blob, meta).Exec(), "Error writing event to Cassandra")
}

const readEvent = `SELECT id, blob, meta FROM events WHERE stream = ? AND id > ? LIMIT ?`

func (sess *CassandraSession) ReadEvents(stream string, cursor string, amount int) ([]Event, error) {
	iter := sess.session.Query(readEvent, stream, cursor, amount).Iter()
	rows, err := iter.SliceMap()
	events := make([]Event, len(rows))
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
