package actions

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type WriteEventFromStream struct {
	Stream string
	Input  io.Reader
}

// Read up to Mb from stream
const MAX_READ = 1024 * 1024

func readJSONFromStream(input io.Reader) ([]byte, error) {
	inputBuf := make([]byte, MAX_READ)
	_, err := input.Read(inputBuf)
	if err != nil {
		return inputBuf, errors.Wrap(err, "Input read failed")
	}
	for i, v := range inputBuf {
		if v == '\x00' {
			inputBuf = inputBuf[:i]
			break
		}
	}
	var js map[string]interface{}
	err = json.Unmarshal(inputBuf, &js)
	if err != nil {
		return inputBuf, errors.Wrap(err, "Input is not JSON")
	}
	out, err := json.Marshal(js)
	if err != nil {
		return inputBuf, errors.Wrap(err, "Marshalling as JSON failed")
	}
	return out, nil
}

func (wes WriteEventFromStream) Run(c Context) error {
	we := WriteEvent{Stream: wes.Stream, Meta: make(map[string]string)}
	we.Meta["origin"] = "stream"
	we.Meta["compressed"] = "false"
	var err error
	we.Payload, err = readJSONFromStream(wes.Input)
	if err != nil {
		return errors.Wrap(err, "Failed reading input")
	}
	return errors.Wrap(we.Run(c), "Failed running WriteEvent action")
}
