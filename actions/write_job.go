package actions

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/simplejson"
)

type WriteJob struct {
	Payload simplejson.Object `json:"payload"`
	Stream  string            `json:"stream"`
}

func (wj WriteJob) Run(c Context) error {
	payload, err := json.Marshal(wj.Payload)
	if err != nil {
		return errors.Wrap(err, "WriteJob jsonMarshal payload")
	}

	if err := persistEvent(wj.Stream, payload, "http; write_job", c); err != nil {
		return errors.Wrap(err, "WriteJob persistEvent")
	}
	c.Telemetry().Incr("write", []string{"stream:" + wj.Stream})
	return nil
}

func (wj WriteJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"payload": wj.Payload, "stream": wj.Stream})
}
