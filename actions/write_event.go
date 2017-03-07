package actions

import (
	"fmt"

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
