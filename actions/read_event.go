package actions

import (
	"fmt"

	"github.com/pkg/errors"
)

type ReadEvent struct {
	Stream string
	Cursor string
	Amount int
}

const minTimeuuid = "00000000-0000-1000-8080-808080808080"

func (re ReadEvent) Run(c Context) error {
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
		js := string(ev.Blob)
		fmt.Printf("%s - %s: {%v} %s\n", ev.Stream, ev.ID, ev.Meta, js)
	}
	return errors.Wrap(err, fmt.Sprintf("Error reading event from stream %s", re.Stream))
}
