package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

type WriteEvent struct {
	Stream  string
	Payload []byte
	Meta    map[string]string
}

func (we WriteEvent) Run(c Context) error {
	session, err := c.Persistence().Session()
	if err != nil {
		return errors.Wrap(err, "Obtaining session failed")
	}
	we.Meta["version"] = c.Version()
	return errors.Wrap(session.WriteEvent(we.Stream, we.Payload, we.Meta),
		fmt.Sprintf("Error writing event to stream %s", we.Stream))
}

func (we WriteEvent) MarshalJSON() ([]byte, error) {
	return []byte(`"Insert event to Cassandra stream"`), nil
}

func persistEvent(stream string, payload []byte, origin string, c Context) error {
	var err error
	for i := 0; i < c.Retries(); i += 1 {
		action := WriteEvent{Stream: stream, Payload: payload, Meta: make(map[string]string)}
		action.Meta["origin"] = origin
		action.Meta["compressed"] = "false"
		action.Meta["time"] = string(time.Now().Unix())
		err = action.Run(c)
		if err == nil {
			break
		}
		time.Sleep(c.Backoff(i))
		c.Telemetry().Incr("persist.retries", []string{"stream" + stream})
		log.Println("Retrying persistEvent, because of", err)
	}
	return err
}
