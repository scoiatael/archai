package actions

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/scoiatael/archai/simplejson"
)

type BulkWriteJob struct {
	Schema  []interface{} `json:"schema"`
	Objects []interface{} `json:"data"`
	Stream  string        `json:"stream"`
}

func (wj BulkWriteJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"schema": wj.Schema,
		"data":   wj.Objects,
		"stream": wj.Stream,
	})
}

func makeObjectWithSchema(obj interface{}, schema []interface{}) (simplejson.Object, error) {
	object_with_schema := make(simplejson.Object)
	object, conv := obj.([]interface{})
	if !conv {
		return object_with_schema, fmt.Errorf("Failed to convert obj to array")
	}
	for j, name := range schema {
		name, conv := name.(string)
		if !conv {
			return object_with_schema, fmt.Errorf("%d: Failed to convert schema value to string", j)
		}
		if len(object) <= j {
			return object_with_schema, fmt.Errorf("%d: Not enough values", j)
		}
		object_with_schema[name] = object[j]
	}
	return object_with_schema, nil
}

func (wj BulkWriteJob) Run(c Context) error {
	c.Telemetry().Incr("bulk_write.aggregate.attempt", []string{"stream:" + wj.Stream})
	for i, obj := range wj.Objects {
		object, err := makeObjectWithSchema(obj, wj.Schema)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("HTTP server splitting payload to bulk_write event at %d", i))
		}
		payload, err := json.Marshal(object)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("HTTP server marshalling payload to bulk_write event at %d", i))
		}
		if err := persistEvent(wj.Stream, payload, "http; bulk_write_job", c); err != nil {
			return errors.Wrap(err, "HTTP server bulk_writing events")
		}
		c.Telemetry().Incr("write", []string{"stream:" + wj.Stream})
	}
	c.Telemetry().Incr("bulk_write.aggregate.write", []string{"stream:" + wj.Stream})
	return nil
}
