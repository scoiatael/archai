package actions

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/simplejson"
)

type WriteEventFromStream struct {
	Stream string
	Input  bufio.Reader
}

func readJSONFromStream(input bufio.Reader) ([]byte, error) {
	buf, err := input.ReadBytes('\n')
	if err != nil {
		return buf, err
	}
	out, err := simplejson.Validate(buf)
	return out, nil
}

func (wes WriteEventFromStream) Run(c Context) error {
	we := WriteEvent{Stream: wes.Stream, Meta: make(map[string]string)}
	we.Meta["origin"] = "stream"
	we.Meta["compressed"] = "false"
	var err error
	for {
		we.Payload, err = readJSONFromStream(wes.Input)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "Failed reading input")
		}
		err = errors.Wrap(we.Run(c), "Failed running WriteEvent action")
		if err != nil {
			return err
		}
	}
}
