package actions

import (
	"fmt"
	"io"

	"github.com/scoiatael/archai/persistence"
)

type ReadEventsToStream struct {
	Stream string
	Cursor string
	Amount int

	Output io.Writer
}

func printEventToStream(out io.Writer, ev persistence.Event) error {
	js := string(ev.Blob)
	str := fmt.Sprintf("%s - %s: {%v} %s\n", ev.Stream, ev.ID, ev.Meta, js)
	_, err := out.Write([]byte(str))
	return err
}

func (res ReadEventsToStream) Run(c Context) error {
	ch := make(chan (persistence.Event), 100)
	err := ReadEvents{Stream: res.Stream, Cursor: res.Cursor, Amount: 10, Output: ch}.Run(c)
	if err != nil {
		return err
	}
	for ev := range ch {
		err := printEventToStream(res.Output, ev)
		if err != nil {
			return err
		}
	}
	return nil
}
