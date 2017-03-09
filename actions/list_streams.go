package actions

import (
	"fmt"

	"github.com/pkg/errors"
)

type ListStreams struct {
}

func (re ListStreams) Run(c Context) error {
	session, err := c.Persistence().Session()
	if err != nil {
		return errors.Wrap(err, "Obtaining session failed")
	}
	streams, err := session.ListStreams()
	for _, s := range streams {
		println(s)
	}
	return errors.Wrap(err, fmt.Sprintf("Error listing streams"))
}

func (re ListStreams) MarshalJSON() ([]byte, error) {
	return []byte(`"List streams"`), nil
}
