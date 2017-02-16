package actions

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/persistence"
)

type ReadEvents struct {
	Stream string
	Cursor string
	Amount int

	Output chan (persistence.Event)
}

const minTimeuuid = "00000000-0000-1000-8080-808080808080"

func (re ReadEvents) Run(c Context) error {
	persistenceProvider := c.Persistence()
	session, err := persistenceProvider.Session()
	if err != nil {
		return errors.Wrap(err, "Obtaining session failed")
	}
	if len(re.Cursor) == 0 {
		re.Cursor = minTimeuuid
	}
	events, err := session.ReadEvents(re.Stream, re.Cursor, re.Amount)
	for _, ev := range events {
		re.Output <- ev
	}
	return errors.Wrap(err, fmt.Sprintf("Error reading event from stream %s", re.Stream))
}
