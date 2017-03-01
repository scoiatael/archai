package actions

import (
	"bufio"

	"github.com/scoiatael/archai/types"
)

type ReadEventsToStream struct {
	Stream string
	Cursor string

	Output bufio.Writer
}

func printEventToStream(out bufio.Writer, ev types.Event) error {
	buf, err := types.EventToJson(ev)
	if err != nil {
		return err
	}
	_, err = out.Write(buf)
	if err != nil {
		return err
	}
	_, err = out.WriteRune('\n')
	if err != nil {
		return err
	}
	err = out.Flush()
	return err
}

func (res ReadEventsToStream) Run(c Context) error {
	size := 10
	cursor := res.Cursor
	for {
		action := ReadEvents{Stream: res.Stream, Cursor: cursor, Amount: size}
		err := action.Run(c)
		events := action.Events
		if err != nil {
			return err
		}
		for _, ev := range events {
			err := printEventToStream(res.Output, ev)
			if err != nil {
				return err
			}
		}
		if len(events) < size {
			break
		}
	}
	return nil
}
