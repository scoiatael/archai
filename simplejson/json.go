package simplejson

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Json interface{}
type Object map[string]Json
type String string
type Array []Json
type ObjectArray []Object

func Read(buf []byte) (Object, error) {
	var js Object

	err := json.Unmarshal(buf, &js)
	return js, errors.Wrap(err, "Input is not JSON")
}

func Write(js Object) ([]byte, error) {
	out, err := json.Marshal(js)
	return out, errors.Wrap(err, "Marshalling as JSON failed")
}

func Validate(buf []byte) ([]byte, error) {
	var (
		err error
	)
	js, err := Read(buf)
	if err != nil {
		return buf, err
	}
	return Write(js)
}
