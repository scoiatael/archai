package actions

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/types"
)

type ReadEvents struct {
	Stream string
	Cursor string
	Amount int

	Events []types.Event
}

const minTimeuuid = "00000000-0000-1000-8080-808080808080"

func (re *ReadEvents) Run(c Context) error {
	persistenceProvider := c.Persistence()
	session, err := persistenceProvider.Session()
	if err != nil {
		return errors.Wrap(err, "Obtaining session failed")
	}
	defer session.Close()
	if len(re.Cursor) == 0 {
		re.Cursor = minTimeuuid
	}
	events, err := session.ReadEvents(re.Stream, re.Cursor, re.Amount)
	re.Events = events
	return errors.Wrap(err, fmt.Sprintf("Error reading event from stream %s", re.Stream))
}
