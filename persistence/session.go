package persistence

import (
	"github.com/gocql/gocql"
)

type Session interface {
	WriteEvent(string, []byte, map[string]string) error
	ReadEvents(string, string, int) []Event
}

type CassandraSession struct {
	session *gocql.Session
}

const insertEvent = `INSERT INTO events (id, stream, blob, meta) VALUES (now(), ?, ?, ?)`

func (sess *CassandraSession) WriteEvent(stream string, blob []byte, meta map[string]string) error {
	return sess.session.Query(insertEvent, stream, blob, meta).Exec()
}

func (sess *CassandraSession) ReadEvents(stream string, cursor string, amount int) []Event {
	events := make([]Event, 0)
	return events
}
