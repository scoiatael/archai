package actions

import (
	"fmt"
	"io"

	"github.com/scoiatael/archai/types"
)

type ReadEventsToStream struct {
	Stream string
	Cursor string

	Output io.Writer
}

func printEventToStream(out io.Writer, ev types.Event) error {
	js := string(ev.Blob)
	str := fmt.Sprintf("%s - %s: {%v} %s\n", ev.Stream, ev.ID, ev.Meta, js)
	_, err := out.Write([]byte(str))
	return err
}

func (res ReadEventsToStream) Run(c Context) error {
	action := ReadEvents{Stream: res.Stream, Cursor: res.Cursor, Amount: 10}
	err := action.Run(c)
	events := action.Events
	if err != nil {
		return err
	}
	res.Output.Write([]byte(fmt.Sprintln("STREAM -- ID -- Meta -- Blob")))
	for _, ev := range events {
		err := printEventToStream(res.Output, ev)
		if err != nil {
			return err
		}
	}
	return nil
}
